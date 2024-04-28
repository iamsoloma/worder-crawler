package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string `yaml:"env" env-default:"local" env-required:"true"`
	UserAgent string `yaml:"userAgent" env-default:"Worder/0.1"`
	Surreal   `yaml:"surreal"`
}

type Surreal struct {
	//UserAgent     string        `yaml:"userAgent" env-default:"Worder/0.1"`
}

func GetConf() *Config {
	configPath := "./metadata/config.yaml"
	if configPath == "" {
		log.Fatalln("Config is not set!")
	}

	//Check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist!", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Can`t read config: %s", configPath)
	}

	return &cfg
}
