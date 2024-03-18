package main

import (
	"fmt"
	"os"

	"github.com/epiefe/jswap/internal/adoptium"
)

func main() {
	//availableReleases()
	downloadRelease(21)
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
	err := adoptium.DownloadRelease(release)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
