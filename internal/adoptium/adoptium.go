package adoptium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/epiefe/jswap/internal/adoptium/system"
	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/util"
	"github.com/epiefe/jswap/internal/web"
)

const api = "https://api.adoptium.net/v3"

func PrintLocalReleases(major int) error {
	config, err := util.ReadJswapConfig()
	if err != nil {
		return err
	}
	if major > 0 {
		config.JDKs = slices.DeleteFunc(config.JDKs, func(jdk util.JDKInfo) bool { return jdk.Major != major })
	}
	for _, jdk := range config.JDKs {
		fmt.Println(jdk.Release)
	}
	if len(config.JDKs) == 0 {
		fmt.Println("            N/A")
	}
	return nil
}

// Prints available releases for a specific major. If major is 0,
// then prints every available release.
func PrintRemoteReleases(major int) error {
	config, err := util.ReadJswapConfig()
	if err != nil {
		return err
	}
	page := 0
	next := true
	for next {
		url := fmt.Sprintf(
			"%s/info/release_names?architecture=%s&image_type=jre&os=%s&page=%d&page_size=20&release_type=ga&sort_method=DEFAULT&sort_order=ASC&vendor=eclipse",
			api, system.ARCH, system.OS, page,
		)
		if major > 0 {
			url += fmt.Sprintf("&version=[%d,%d)", major, major+1)
		}
		result, headers, err := web.FetchJson[releases](url)
		if err != nil {
			return err
		}
		for _, release := range result.Releases {
			if slices.ContainsFunc(config.JDKs, func(jdk util.JDKInfo) bool { return jdk.Release == release }) {
				release += " [installed]"
			}
			fmt.Println(release)
		}
		next = strings.Contains(headers.Get("link"), "rel=\"next\"")
		page += 1
	}
	if page == 0 {
		fmt.Println("            N/A")
	}
	return nil
}

// Downloads the latest JDK release for a specific major
func DownloadLatestRelease(major int) error {
	// Get latest release info from api
	fmt.Printf("Searching latest JDK %d for %s %s\n", major, system.OS, system.ARCH)
	url := fmt.Sprintf("%s/assets/latest/%d/hotspot?architecture=%s&image_type=jre&os=%s&vendor=eclipse", api, major, system.ARCH, system.OS)
	assets, _, err := web.FetchJson[[]asset](url)
	if err != nil {
		return err
	}
	if len(*assets) == 0 {
		return fmt.Errorf("no assets available for major %d", major)
	}
	asset := (*assets)[0]
	fmt.Printf("Found release %s\n", asset.ReleaseName)

	// Download and install
	path, err := getFromLink(asset.Binary.Package.Link)
	if err != nil {
		return err
	}
	// Update jswap.json file
	if err = util.StoreJDKConfig(util.JDKInfo{
		Vendor:      "adoptium",
		Major:       asset.Version.Major,
		Release:     asset.ReleaseName,
		ReleaseDate: asset.Binary.UpdatedAt,
		Path:        path,
	}); err != nil {
		return err
	}
	fmt.Printf("Successfully installed %s\n", asset.ReleaseName)
	return nil
}

// Downloads a specific JDK release
func DownloadRelease(name string) error {
	// Get release info from api
	url := fmt.Sprintf("%s/assets/release_name/eclipse/%s?architecture=%s&image_type=jre&os=%s", api, name, system.ARCH, system.OS)
	release, _, err := web.FetchJson[release](url)
	if err != nil {
		return err
	}
	if len(release.Binaries) == 0 {
		return fmt.Errorf("no binaries available for %s", name)
	}
	// Download and install
	path, err := getFromLink(release.Binaries[0].Package.Link)
	if err != nil {
		return err
	}
	// Update jswap.json file
	if err = util.StoreJDKConfig(util.JDKInfo{
		Vendor:      "adoptium",
		Major:       release.VersionData.Major,
		Release:     name,
		ReleaseDate: release.Binaries[0].UpdatedAt,
		Path:        path,
	}); err != nil {
		return err
	}
	fmt.Printf("Successfully installed %s\n", name)
	return nil
}

func getFromLink(link string) (string, error) {
	// Download release archive
	cacheDir := file.CacheDir()
	defer os.RemoveAll(cacheDir)
	archive, err := web.DownloadFile(link, filepath.Join(cacheDir, "archive"))
	if err != nil {
		return "", err
	}

	// Extract archive
	fmt.Println("Extracting archive...")
	extractDir := filepath.Join(cacheDir, "extracted")
	if err := file.ExtractArchive(archive, extractDir); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return "", err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return "", errors.New("unexpected archive structure")
	}
	name := entries[0].Name()
	extractedPath := filepath.Join(extractDir, name)

	// Create adoptium folder if it does not exist
	jdkDir := filepath.Join(file.JswapHome(), "jdks", "adoptium")
	if err = os.MkdirAll(jdkDir, os.ModePerm); err != nil {
		return "", err
	}

	// Eventually remove pre-existing folder in same path
	jdkPath := filepath.Join(jdkDir, name)
	if err = os.RemoveAll(jdkPath); err != nil {
		return "", err
	}

	// Move extracted jdk to adoptium folder
	if err = os.Rename(extractedPath, jdkPath); err != nil {
		return "", err
	}
	return jdkPath, nil
}
