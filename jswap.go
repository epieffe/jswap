package main

import (
	"fmt"
	"os"

	"github.com/epiefe/jswap/internal/adoptium"
)

func main() {
	//availableReleases()
	//downloadRelease(21)
	downloadVersion("jdk8u402-b06")
}

func availableReleases() {
	releases, err := adoptium.AvailableReleases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	for _, release := range releases {
		fmt.Println(release)
	}
}

func downloadRelease(release int) {
	if err := adoptium.DownloadLatest(release); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}

func downloadVersion(version string) {
	if err := adoptium.DownloadVersion(version); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
