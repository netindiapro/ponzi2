// Package config provides functions for loading and saving the user's configuration.
package config

import (
	"encoding/gob"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/golang/glog"
)

// Config configures the app.
type Config struct {
	CurrentStock *Stock
	Stocks       []*Stock
}

// Stock identifies a single stock by symbol.
type Stock struct {
	Symbol string
}

// Load loads the user's config from disk.
func Load() (*Config, error) {
	cfgPath, err := userConfigPath()
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("loading config from %s", cfgPath)

	file, err := os.Open(cfgPath)
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	dec := gob.NewDecoder(file)
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save saves the user's config to disk.
func Save(cfg *Config) error {
	cfgPath, err := userConfigPath()
	if err != nil {
		return err
	}

	glog.V(2).Infof("saving config to %s", cfgPath)

	file, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	return nil
}

func userConfigPath() (string, error) {
	dirPath, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dirPath, "config.gob"), nil
}

func userConfigDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	p := filepath.Join(u.HomeDir, ".config", "ponzi")
	if err := os.MkdirAll(p, 0755); err != nil {
		return "", err
	}
	return p, nil
}
