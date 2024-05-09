package jdk

import (
	"fmt"
	"os"
	"slices"

	"github.com/epiefe/jswap/internal/file"
	"github.com/epiefe/jswap/internal/jdk/adoptium"
)

// ListLocal prints all the installed JDK releases, matching a given major.
// If major is 0 prints all the installed JDKs.
func ListLocal(major int) error {
	conf, err := readJswapConfig()
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
	conf, err := readJswapConfig()
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
		if slices.ContainsFunc(conf.JDKs, func(info *JDKInfo) bool { return info.Release == release }) {
			release += " [installed]"
		}
		fmt.Println(release)
	}
	return nil
}

// GetLatest downloads and installs the latest JDK
// release matching a given major
func GetLatest(major int) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	asset, err := adoptium.GetLatestAsset(major)
	if err != nil {
		return err
	}
	fmt.Printf("Found release %s\n", asset.ReleaseName)
	// Check if release is already installed
	if checkConflict(asset.ReleaseName, conf) {
		fmt.Printf("Latest JDK %d is already installed\n", major)
		return nil
	}
	path, err := adoptium.GetFromLink(asset.Binary.Package.Link)
	if err != nil {
		return err
	}
	return installJDK(&JDKInfo{
		Vendor:      "adoptium",
		Major:       asset.Version.Major,
		Release:     asset.ReleaseName,
		ReleaseDate: asset.Binary.UpdatedAt,
		Path:        path,
	})
}

// GetRelease downloads and installs a specific JDK release.
func GetRelease(name string) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	// Check if release is already installed
	if checkConflict(name, conf) {
		fmt.Printf("release %s is already installed\n", name)
		return nil
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
	return installJDK(&JDKInfo{
		Vendor:      "adoptium",
		Major:       release.VersionData.Major,
		Release:     name,
		ReleaseDate: release.Binaries[0].UpdatedAt,
		Path:        path,
	})
}

// SetMajor sets the latest installed release matching
// a given major as the current JDK.
func SetMajor(major int) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	var info *JDKInfo
	for _, v := range conf.JDKs {
		if v.Major == major && (info == nil || v.ReleaseDate > info.ReleaseDate) {
			info = v
		}
	}
	if info == nil {
		return fmt.Errorf("no installed release found for JDK %d", major)
	}
	return setJDK(info, conf)
}

// SetRelease sets a specific JDK release as the current JDK.
func SetRelease(name string) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	var info *JDKInfo
	for _, v := range conf.JDKs {
		if v.Release == name {
			info = v
			break
		}
	}
	if info == nil {
		return fmt.Errorf("release %s is not installed", name)
	}
	return setJDK(info, conf)
}

func RemoveReleases(names ...string) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	for _, name := range names {
		i := slices.IndexFunc(conf.JDKs, func(v *JDKInfo) bool { return v.Release == name })
		if i == -1 {
			fmt.Fprintf(os.Stderr, "Warning: release %s not found\n", name)
			continue
		}
		if err := os.RemoveAll(conf.JDKs[i].Path); err != nil {
			return err
		}
		if conf.CurrentJDK != nil && conf.JDKs[i].Path == conf.CurrentJDK.Path {
			conf.CurrentJDK = nil
			if err := os.RemoveAll(file.JavaHome()); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s", err)
			}
		}
		conf.JDKs = slices.Delete(conf.JDKs, i, i+1)
		writeJswapConfig(conf)
		fmt.Printf("Removed %s\n", name)
	}
	return nil
}

func setJDK(info *JDKInfo, conf *JswapConfig) error {
	if err := file.Link(info.Path, file.JavaHome()); err != nil {
		return err
	}
	fmt.Printf("Now using JDK %d (release %s)\n", info.Major, info.Release)
	conf.CurrentJDK = info
	writeJswapConfig(conf)
	return nil
}

func installJDK(info *JDKInfo) error {
	conf, err := readJswapConfig()
	if err != nil {
		return err
	}
	conf.JDKs = append(conf.JDKs, info)
	if err = writeJswapConfig(conf); err != nil {
		return err
	}
	fmt.Printf("Successfully installed %s\n", info.Release)
	if conf.CurrentJDK == nil {
		// Since current JDK is not set, we use this one
		setJDK(info, conf)
	}
	return nil
}

func checkConflict(release string, conf *JswapConfig) bool {
	return slices.ContainsFunc(conf.JDKs, func(info *JDKInfo) bool {
		return info.Release == release
	})
}
