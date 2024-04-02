package adoptium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/jdk/adoptium/system"
	"github.com/epiefe/jswap/internal/web"
)

const api = "https://api.adoptium.net/v3"

// ReleaseNames returns all the available releases for a given major.
// If major is 0 retrieves every available release.
func ReleaseNames(major int) ([]string, error) {
	page := 0
	next := true
	result := []string{}
	for next {
		url := fmt.Sprintf(
			"%s/info/release_names?architecture=%s&image_type=jre&os=%s&page=%d&page_size=20&release_type=ga&sort_method=DEFAULT&sort_order=ASC&vendor=eclipse",
			api, system.ARCH, system.OS, page,
		)
		if major > 0 {
			url += fmt.Sprintf("&version=[%d,%d)", major, major+1)
		}
		fetched, headers, err := web.FetchJson[releases](url)
		if err != nil {
			return nil, err
		}
		result = append(result, fetched.Releases...)
		next = strings.Contains(headers.Get("link"), "rel=\"next\"")
		page += 1
	}
	return result, nil
}

// GetLatestAsset returns informations about the latest JDK
// release for a specific major.
func GetLatestAsset(major int) (*Asset, error) {
	// Get latest release info from api
	fmt.Printf("Searching latest JDK %d for %s %s\n", major, system.OS, system.ARCH)
	url := fmt.Sprintf("%s/assets/latest/%d/hotspot?architecture=%s&image_type=jre&os=%s&vendor=eclipse", api, major, system.ARCH, system.OS)
	assets, _, err := web.FetchJson[[]Asset](url)
	if err != nil {
		return nil, err
	}
	if len(*assets) == 0 {
		return nil, fmt.Errorf("no assets available for major %d", major)
	}
	return &(*assets)[0], nil
}

// GetRelease returns information about a specific JDK release.
func GetRelease(name string) (*Release, error) {
	// Get release info from api
	url := fmt.Sprintf("%s/assets/release_name/eclipse/%s?architecture=%s&image_type=jre&os=%s", api, name, system.ARCH, system.OS)
	release, _, err := web.FetchJson[Release](url)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func GetFromLink(link string) (string, error) {
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
