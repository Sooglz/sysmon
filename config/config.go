package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval   int    `yaml:"interval"`
	ListenAddr string `yaml:"listen_addr"`
	DiskPath   string `yaml:"disk_path"`
}

func Default() Config {
	return Config{
		Interval:   5,
		ListenAddr: "127.0.0.1:9100",
		DiskPath:   "/",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be > 0")
	}
	if c.ListenAddr == "" {
		return fmt.Errorf("listen_addr must not be empty")
	}
	if c.DiskPath == "" {
		return fmt.Errorf("disk_path must not be empty")
	}
	return nil
}
