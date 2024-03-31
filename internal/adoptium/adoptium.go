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

// Fetches available Java releases from Eclipse Temurin.
func AvailableReleases() ([]string, error) {
	url := api + "/info/available_releases"
	available, _, err := web.FetchJson[availableReleases](url)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(available.Releases))
	for i, v := range available.Releases {
		if slices.Contains(available.LTS, v) {
			result[i] = fmt.Sprintf("Adoptium JDK %d - LTS", v)
		} else {
			result[i] = fmt.Sprintf("Adoptium JDK %d", v)
		}
	}
	return result, nil
}

// Prints available versions for a specific release. If release is 0,
// then prints available versions for every release.
func PrintRemoteVersions(release int) error {
	page := 0
	next := true
	for next {
		url := fmt.Sprintf(
			"%s/info/release_names?architecture=%s&image_type=jre&os=%s&page=%d&page_size=20&release_type=ga&sort_method=DEFAULT&sort_order=ASC&vendor=eclipse",
			api, system.ARCH, system.OS, page,
		)
		if release > 0 {
			url += fmt.Sprintf("&version=[%d,%d)", release, release+1)
		}
		result, headers, err := web.FetchJson[versions](url)
		if err != nil {
			return err
		}
		for _, version := range result.Releases {
			fmt.Println(version)
		}
		next = strings.Contains(headers.Get("link"), "rel=\"next\"")
		page += 1
	}
	if page == 0 {
		fmt.Println("            N/A")
	}
	return nil
}

// Downloads the latest JDK build for a specific release
func DownloadLatest(release int) error {
	// Get latest release info from api
	fmt.Printf("Searching latest JDK %d for %s %s\n", release, system.OS, system.ARCH)
	url := fmt.Sprintf("%s/assets/latest/%d/hotspot?architecture=%s&image_type=jre&os=%s&vendor=eclipse", api, release, system.ARCH, system.OS)
	assets, _, err := web.FetchJson[[]asset](url)
	if err != nil {
		return err
	}
	if len(*assets) == 0 {
		return fmt.Errorf("no assets available for release %d", release)
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
	return nil
}

// Downloads a specific JDK version
func DownloadVersion(version string) error {
	// Get release version info from api
	url := fmt.Sprintf("%s/assets/release_name/eclipse/%s?architecture=%s&image_type=jre&os=%s", api, version, system.ARCH, system.OS)
	release, _, err := web.FetchJson[release](url)
	if err != nil {
		return err
	}
	if len(release.Binaries) == 0 {
		return fmt.Errorf("no binaries available for %s", version)
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
		Release:     version,
		ReleaseDate: release.Binaries[0].UpdatedAt,
		Path:        path,
	}); err != nil {
		return err
	}
	return nil
}

func getFromLink(link string) (string, error) {
	// Download latest release archive
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
	jdkDir := filepath.Join(file.JswapDir(), "jdk", "adoptium")
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
	fmt.Printf("Successfully installed %s\n", name)
	return jdkPath, nil
}
