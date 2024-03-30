package adoptium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/epiefe/jswap/internal/adoptium/system"
	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/web"
)

const api = "https://api.adoptium.net/v3"

// Fetches available Java releases from Eclipse Temurin.
func AvailableReleases() ([]string, error) {
	url := api + "/info/available_releases"
	available, err := web.FetchJson[availableReleases](url)
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

// Downloads the latest JDK build for a specific release
func DownloadLatest(release int) error {
	// Get latest release info from api
	fmt.Printf("Searching latest JDK %d for %s %s\n", release, system.OS, system.ARCH)
	url := fmt.Sprintf("%s/assets/latest/%d/hotspot?architecture=%s&image_type=jre&os=%s&vendor=eclipse", api, release, system.ARCH, system.OS)
	assets, err := web.FetchJson[[]asset](url)
	if err != nil {
		return err
	}
	if len(*assets) == 0 {
		return fmt.Errorf("no assets available for release %d", release)
	}
	asset := (*assets)[0]
	fmt.Printf("Found release %s\n", asset.ReleaseName)

	// Download and install
	if err := getFromLink(asset.Binary.Package.Link); err != nil {
		return err
	}
	return nil
}

// Downloads a specific JDK version
func DownloadVersion(version string) error {
	// Get release version info from api
	url := fmt.Sprintf("%s/assets/release_name/eclipse/%s?architecture=%s&image_type=jre&os=%s", api, version, system.ARCH, system.OS)
	release, err := web.FetchJson[release](url)
	if err != nil {
		return err
	}
	if len(release.Binaries) == 0 {
		return fmt.Errorf("no binaries available for %s", version)
	}
	// Download and install
	if err := getFromLink(release.Binaries[0].Package.Link); err != nil {
		return err
	}
	return nil
}

func getFromLink(link string) error {
	// Download latest release archive
	cacheDir := file.CacheDir()
	defer os.RemoveAll(cacheDir)
	archive, err := web.DownloadFile(link, filepath.Join(cacheDir, "archive"))
	if err != nil {
		return err
	}

	// Extract archive
	fmt.Println("Extracting archive...")
	extractDir := filepath.Join(cacheDir, "extracted")
	if err := file.ExtractArchive(archive, extractDir); err != nil {
		return err
	}
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return errors.New("unexpected archive structure")
	}
	name := entries[0].Name()
	extractedPath := filepath.Join(extractDir, name)

	// Create adoptium folder if it does not exists
	jdkDir := filepath.Join(file.JswapDir(), "jdk", "adoptium")
	if err = os.MkdirAll(jdkDir, os.ModePerm); err != nil {
		return err
	}

	// Eventually remove pre-existing jdk with same version
	jdkPath := filepath.Join(jdkDir, name)
	if err = os.RemoveAll(jdkPath); err != nil {
		return err
	}

	// Move extracted jdk to adoptium folder
	if err = os.Rename(extractedPath, jdkPath); err != nil {
		return err
	}
	fmt.Printf("Successfully installed %s\n", name)
	return nil
}
