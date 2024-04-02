package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func JavaHome() string {
	return filepath.Join(JswapHome(), "current-jdk")
}

func JswapHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
		os.Exit(1)
	}
	return filepath.Join(dir, ".jswap")
}

func CacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = filepath.Join(JswapHome(), ".cache")
	}
	return filepath.Join(dir, "jswap")
}
