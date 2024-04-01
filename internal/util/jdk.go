package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/epiefe/jswap/internal/file"
)

func UseMajor(major int) error {
	config, err := ReadJswapConfig()
	if err != nil {
		return err
	}
	var jdk *JDKInfo
	for _, v := range config.JDKs {
		if v.Major == major && (jdk == nil || v.ReleaseDate > jdk.ReleaseDate) {
			jdk = &v
		}
	}
	if jdk == nil {
		return fmt.Errorf("no installed release found for JDK %d", major)
	}
	return useJDK(jdk)
}

func UseRelease(name string) error {
	config, err := ReadJswapConfig()
	if err != nil {
		return err
	}
	var jdk *JDKInfo
	for _, v := range config.JDKs {
		if v.Release == name {
			jdk = &v
			break
		}
	}
	if jdk == nil {
		return fmt.Errorf("release %s is not installed", name)
	}
	return useJDK(jdk)
}

func useJDK(jdk *JDKInfo) error {
	if err := os.MkdirAll(file.CacheDir(), os.ModePerm); err != nil {
		return err
	}
	// Update JAVA_HOME symlink in the most possible atomic way
	// https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go
	symlinkPathTmp := filepath.Join(file.CacheDir(), "jdklink.tmp")
	if err := os.Remove(symlinkPathTmp); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(jdk.Path, symlinkPathTmp); err != nil {
		return err
	}

	if err := os.Rename(symlinkPathTmp, file.JavaHome()); err != nil {
		return err
	}
	fmt.Printf("Now using JDK %d (release %s)\n", jdk.Major, jdk.Release)
	return nil
}
