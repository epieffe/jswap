package jdk

import (
	"fmt"

	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/jdk/adoptium"
	"github.com/epiefe/jswap/internal/jdk/config"
)

func GetLatest(major int) error {
	info, err := adoptium.DownloadLatestRelease(major)
	if err != nil {
		return err
	}
	return installJDK(info)
}

func GetRelease(release string) error {
	info, err := adoptium.DownloadRelease(release)
	if err != nil {
		return err
	}
	return installJDK(info)
}

func UseMajor(major int) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	var info *config.JDKInfo
	for _, v := range conf.JDKs {
		if v.Major == major && (info == nil || v.ReleaseDate > info.ReleaseDate) {
			info = v
		}
	}
	if info == nil {
		return fmt.Errorf("no installed release found for JDK %d", major)
	}
	return useJDK(info, conf)
}

func UseRelease(name string) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	var info *config.JDKInfo
	for _, v := range conf.JDKs {
		if v.Release == name {
			info = v
			break
		}
	}
	if info == nil {
		return fmt.Errorf("release %s is not installed", name)
	}
	return useJDK(info, conf)
}

func useJDK(info *config.JDKInfo, conf *config.JswapConfig) error {
	if err := file.Link(info.Path, file.JavaHome()); err != nil {
		return err
	}
	fmt.Printf("Now using JDK %d (release %s)\n", info.Major, info.Release)
	conf.CurrentJDK = info
	config.WriteJswapConfig(conf)
	return nil
}

func installJDK(info *config.JDKInfo) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	conf.AddJDK(info)
	if err = config.WriteJswapConfig(conf); err != nil {
		return err
	}
	fmt.Printf("Successfully installed %s\n", info.Release)
	if conf.CurrentJDK == nil {
		// Since current JDK is not set, we use this one
		useJDK(info, conf)
	}
	return nil
}
