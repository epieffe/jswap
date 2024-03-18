package adoptium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"

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
func DownloadRelease(version int) error {
	url := fmt.Sprintf(
		"%s/assets/latest/%d/hotspot?architecture=x64&image_type=jre&os=%s&vendor=eclipse",
		api, version, runtime.GOOS,
	)
	// Get latest release info from api
	fmt.Printf("Searching latest JDK %d release for %s %s\n", version, runtime.GOOS, runtime.GOARCH)
	assets, err := web.FetchJson[[]asset](url)
	if err != nil {
		return err
	}
	if len(*assets) == 0 {
		return fmt.Errorf("no assets available for release %d", version)
	}
	asset := (*assets)[0]
	fmt.Printf("Found release %s\n", asset.ReleaseName)

	// Download latest release archive
	cacheDir := file.CacheDir()
	defer os.RemoveAll(cacheDir)
	downloadPath := filepath.Join(cacheDir, "archive", asset.Binary.Package.Name)
	err = web.DownloadFile(asset.Binary.Package.Link, downloadPath)
	if err != nil {
		return err
	}

	// Extract archive
	fmt.Println("Extracting archive...")
	extractDir := filepath.Join(cacheDir, "extracted")
	err = file.ExtractArchive(downloadPath, extractDir)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return errors.New("unexpected archive structure")
	}
	jdkPath := filepath.Join(extractDir, entries[0].Name())

	// Move extracted jdk to jswap folder
	dir := filepath.Join(file.JswapDir(), "jdk", "adoptium")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Rename(jdkPath, filepath.Join(dir, strconv.Itoa(version)))
	if err != nil {
		return err
	}

	fmt.Printf("JDK %d installed successfully!\n", version)
	return nil
}
