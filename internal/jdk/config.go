package jdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/epiefe/jswap/internal/file"
)

type JswapConfig struct {
	JDKs []JDKInfo `json:"jdks"`
}

type JDKInfo struct {
	Vendor      string `json:"vendor"`
	Major       int    `json:"major"`
	Release     string `json:"release"`
	ReleaseDate string `json:"releaseDate"`
	Path        string `json:"path"`
}

func (config *JswapConfig) AddJDK(info JDKInfo) {
	conflict := slices.ContainsFunc(config.JDKs, func(e JDKInfo) bool { return e.Release == info.Release })
	if conflict {
		fmt.Fprintf(os.Stderr, "Warning: release %s was already installed and this will override previous configuration\n", info.Release)
		config.JDKs = slices.DeleteFunc(config.JDKs, func(e JDKInfo) bool { return e.Release == info.Release })
	}
	config.JDKs = append(config.JDKs, info)
}

func (config *JswapConfig) RemoveJDK(release string) {
	oldLen := len(config.JDKs)
	config.JDKs = slices.DeleteFunc(config.JDKs, func(e JDKInfo) bool { return e.Release == release })
	if oldLen == len(config.JDKs) {
		fmt.Fprintf(os.Stderr, "Warning: attempted to remove release %s but was not found\n", release)
	}
}

func defaultConfig() *JswapConfig {
	return &JswapConfig{JDKs: []JDKInfo{}}
}

func StoreJDKConfig(jdk JDKInfo) error {
	config, err := ReadJswapConfig()
	if err != nil {
		return err
	}
	config.AddJDK(jdk)
	if err = WriteJswapConfig(config); err != nil {
		return err
	}
	return nil
}

func ReadJswapConfig() (*JswapConfig, error) {
	file, err := os.Open(filepath.Join(file.JswapHome(), "jswap.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig(), nil
		} else {
			return nil, err
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config JswapConfig
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func WriteJswapConfig(config *JswapConfig) error {
	file, err := os.Create(filepath.Join(file.JswapHome(), "jswap.json"))
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&config)
	if err != nil {
		return err
	}
	return nil
}
