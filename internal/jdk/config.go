package jdk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/epiefe/jswap/internal/file"
)

type JswapConfig struct {
	CurrentJDK *JDKInfo   `json:"currentJDK"`
	JDKs       []*JDKInfo `json:"jdks"`
}

type JDKInfo struct {
	Vendor      string `json:"vendor"`
	Major       int    `json:"major"`
	Release     string `json:"release"`
	ReleaseDate string `json:"releaseDate"`
	Path        string `json:"path"`
}

func readJswapConfig() (*JswapConfig, error) {
	file, err := os.Open(filepath.Join(file.JswapData(), "jswap.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &JswapConfig{JDKs: []*JDKInfo{}}, nil
		} else {
			return nil, err
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var conf JswapConfig
	if err = decoder.Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func writeJswapConfig(conf *JswapConfig) error {
	file, err := os.Create(filepath.Join(file.JswapData(), "jswap.json"))
	if err != nil {
		return err
	}
	defer file.Close()
	// Sort jdks
	sort.Slice(conf.JDKs, func(i, j int) bool {
		return conf.JDKs[i].Major < conf.JDKs[j].Major ||
			(conf.JDKs[i].Major == conf.JDKs[j].Major &&
				conf.JDKs[i].ReleaseDate < conf.JDKs[j].ReleaseDate)
	})
	// Save conf to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&conf)
	if err != nil {
		return err
	}
	return nil
}
