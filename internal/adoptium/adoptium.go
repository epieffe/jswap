package adoptium

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

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
	assets, err := web.FetchJson[[]asset](url)
	if err != nil {
		return err
	}
	if len(*assets) == 0 {
		return fmt.Errorf("no assets available for release %d", version)
	}
	asset := (*assets)[0]

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".jswap", "adoptium", asset.Binary.Package.Name)
	err = web.DownloadFile(asset.Binary.Package.Link, path)
	if err != nil {
		return err
	}
	return nil
}
