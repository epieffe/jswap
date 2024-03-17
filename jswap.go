package main

import (
	"fmt"

	"github.com/epiefe/jswap/internal/adoptium"
)

func main() {
	downloadRelease(21)
}

func availableReleases() {
	releases, err := adoptium.AvailableReleases()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, release := range releases {
		fmt.Println(release)
	}
}

func downloadRelease(release int) {
	err := adoptium.DownloadRelease(release)
	if err != nil {
		fmt.Println(err)
		return
	}
}
