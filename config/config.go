// Package config 负责 SysMon 配置文件的解析、默认值与校验。
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示 SysMon 的运行配置。
type Config struct {
	Interval   int    `yaml:"interval"`
	ListenAddr string `yaml:"listen_addr"`
	DiskPath   string `yaml:"disk_path"`
}

// Default 返回一份带有默认值的配置。
func Default() Config {
	return Config{
		Interval:   5,
		ListenAddr: "127.0.0.1:9100",
		DiskPath:   "/",
	}
}

// Load 从指定路径读取并解析 YAML 配置文件。
// 配置文件路径由管理员通过 --config 参数指定，属本地可信输入。
func Load(path string) (Config, error) {
	cfg := Default()
	// #nosec G304 -- path 来自命令行 --config，由管理员控制，非用户输入
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	return cfg, nil
}

// Validate 校验配置项的合法性。
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
