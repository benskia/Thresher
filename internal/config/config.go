package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// Description:
//	Manages information about config files and existing profiles.
//
// Responsibilities:
//	- Create config file if it doesn't exist.
//	- Read profiles from existing config file.
//	- Write profiles to config file.
//	- Store Profiles during program execution.

type Profile struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Profiles map[string]Profile

type Config struct {
	configPath string
	Profiles   map[string]Profile `json:"profiles"`
}

var Defaults Profiles = Profiles{
	"mid": {
		Name:  "mid",
		Start: 40,
		End:   50,
	},
	"high": {
		Name:  "high",
		Start: 70,
		End:   80,
	},
}

// Get a config pointer either using an existing config file or Defaults.
func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{configPath: configPath}

	// We can still use Defaults if we fail to get Profiles from a config file.
	// In case we were expecting a config file, we should return this error.
	err := cfg.readConfigFile()
	if err != nil {
		cfg.Profiles = Defaults
		err = fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, err
}

func (cfg *Config) SaveConfig() error {
	if err := cfg.writeConfigFile(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func (cfg *Config) readConfigFile() error {
	b, err := os.ReadFile(cfg.configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var profiles Profiles
	if err = json.Unmarshal(b, &profiles); err != nil {
		return fmt.Errorf("error unmarshaling: %w", err)
	}

	cfg.Profiles = profiles
	return nil
}

func (cfg *Config) writeConfigFile() error {
	// To work with a config, we'll need the file and directory where it lives.
	if err := os.MkdirAll(path.Dir(cfg.configPath), 0755); err != nil {
		return fmt.Errorf("error making config directories: %w", err)
	}

	f, err := os.Create(cfg.configPath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}

	b, err := json.Marshal(cfg.Profiles)
	if err != nil {
		return fmt.Errorf("error marshaling: %w", err)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
