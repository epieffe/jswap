package jdk

import (
	"fmt"

	"github.com/epiefe/jswap/internal/file"
)

func UseMajor(major int) error {
	config, err := ReadJswapConfig()
	if err != nil {
		return err
	}
	var info *JDKInfo
	for _, v := range config.JDKs {
		if v.Major == major && (info == nil || v.ReleaseDate > info.ReleaseDate) {
			info = &v
		}
	}
	if info == nil {
		return fmt.Errorf("no installed release found for JDK %d", major)
	}
	return useJDK(info)
}

func UseRelease(name string) error {
	config, err := ReadJswapConfig()
	if err != nil {
		return err
	}
	var info *JDKInfo
	for _, v := range config.JDKs {
		if v.Release == name {
			info = &v
			break
		}
	}
	if info == nil {
		return fmt.Errorf("release %s is not installed", name)
	}
	return useJDK(info)
}

func useJDK(info *JDKInfo) error {
	if err := file.Link(info.Path, file.JavaHome()); err != nil {
		return err
	}
	fmt.Printf("Now using JDK %d (release %s)\n", info.Major, info.Release)
	return nil
}
