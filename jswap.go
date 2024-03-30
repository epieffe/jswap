package main

import (
	"fmt"
	"os"

	"github.com/epiefe/jswap/internal/adoptium"
)

func main() {
	printVersions(33)
	//downloadRelease(21)
	//downloadVersion("jdk8u402-b06")
}

func printVersions(release int) {
	if err := adoptium.PrintRemoteVersions(release); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
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
