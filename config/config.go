package config

import (
	"github.com/qibobo/webgo-gin/db"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}
type ServerConfig struct {
	Port int `yaml:"port"`
}
type DBConfig struct {
	DemoDB db.DatabaseConfig `yaml:"demo_db"`
}

func LoadConfig(bytes []byte) (*Config, error) {
	conf := &Config{
		Server: ServerConfig{
			Port: 8080,
		},
	}
	err := yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *Config) Validate() error {
	return nil
}
