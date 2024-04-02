package util

import (
	"fmt"

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
	if err := file.Link(jdk.Path, file.JavaHome()); err != nil {
		return err
	}
	fmt.Printf("Now using JDK %d (release %s)\n", jdk.Major, jdk.Release)
	return nil
}
