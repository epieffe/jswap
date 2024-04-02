package jdk

import (
	"fmt"
	"slices"

	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/jdk/adoptium"
	"github.com/epiefe/jswap/internal/jdk/config"
)

// ListLocal prints all the installed JDK releases, matching a given major.
// If major is 0 prints all the installed JDKs.
func ListLocal(major int) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	none := true
	for _, info := range conf.JDKs {
		if major == 0 || info.Major == major {
			none = false
			if conf.CurrentJDK != nil && info.Release == conf.CurrentJDK.Release {
				fmt.Printf("%s [current]\n", info.Release)
			} else {
				fmt.Println(info.Release)
			}
		}
	}
	if none {
		fmt.Println("            N/A")
	}
	return nil
}

// ListRemote prints all the JDK releases available to install, matching a
// given major. If major is zero prints all available releases.
func ListRemote(major int) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	releases, err := adoptium.ReleaseNames(major)
	if err != nil {
		return err
	}
	if len(releases) == 0 {
		fmt.Println("            N/A")
		return nil
	}
	for _, release := range releases {
		if slices.ContainsFunc(conf.JDKs, func(info *config.JDKInfo) bool { return info.Release == release }) {
			release += " [installed]"
		}
		fmt.Println(release)
	}
	return nil
}

// GetLatest downloads and installs the latest JDK
// release matching a given major
func GetLatest(major int) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	asset, err := adoptium.GetLatestAsset(major)
	if err != nil {
		return err
	}
	fmt.Printf("Found release %s\n", asset.ReleaseName)
	// Check if release is already installed
	if err := checkConflict(asset.ReleaseName, conf, false); err != nil {
		return err
	}
	path, err := adoptium.GetFromLink(asset.Binary.Package.Link)
	if err != nil {
		return err
	}
	return installJDK(&config.JDKInfo{
		Vendor:      "adoptium",
		Major:       asset.Version.Major,
		Release:     asset.ReleaseName,
		ReleaseDate: asset.Binary.UpdatedAt,
		Path:        path,
	})
}

// GetRelease downloads and installs a specific JDK release.
func GetRelease(name string) error {
	conf, err := config.ReadJswapConfig()
	if err != nil {
		return err
	}
	// Check if release is already installed
	if err := checkConflict(name, conf, false); err != nil {
		return err
	}
	release, err := adoptium.GetRelease(name)
	if err != nil {
		return err
	}
	if len(release.Binaries) == 0 {
		return fmt.Errorf("no binaries available for %s", name)
	}
	path, err := adoptium.GetFromLink(release.Binaries[0].Package.Link)
	if err != nil {
		return err
	}
	return installJDK(&config.JDKInfo{
		Vendor:      "adoptium",
		Major:       release.VersionData.Major,
		Release:     name,
		ReleaseDate: release.Binaries[0].UpdatedAt,
		Path:        path,
	})
}

// UseMajor sets the latest installed release matching
// a given major as the current JDK.
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

// UseRelease sets a specific JDK release as the current JDK.
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

func checkConflict(release string, conf *config.JswapConfig, force bool) error {
	if !force && slices.ContainsFunc(conf.JDKs, func(info *config.JDKInfo) bool { return info.Release == release }) {
		return fmt.Errorf("release %s is already installed", release)
	}
	return nil
}
