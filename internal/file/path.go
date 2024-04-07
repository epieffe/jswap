package file

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func JavaHome() string {
	return filepath.Join(JswapData(), "current-jdk")
}

func TempDir() string {
	return filepath.Join(JswapData(), "tmp")
}

func JswapData() string {
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LocalAppData")
		if localAppData == "" {
			fmt.Fprintln(os.Stderr, "Fatal error: %LocalAppData% is not defined")
			os.Exit(1)
		}
		return filepath.Join(localAppData, "Jswap")

	default:
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
			os.Exit(1)
		}
		return filepath.Join(dir, ".jswap")
	}
}
