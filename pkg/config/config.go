package config

import (
	"log"
	"os"
	"sync"

	//external
	"gopkg.in/yaml.v3"
)

type Config struct {
	Transport TransportConfig `yaml:"transport"`
	Database  DatabaseConfig  `yaml:"database"`
}

type TransportConfig struct {
	Rest RESTTransportConfig `yaml:"rest"`
}

type RESTTransportConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Engine                       string `yaml:"engine"`
	URL                          string `yaml:"url"`
	MaxIdleConnection            int    `yaml:"max_idle_connection"`
	MaxOpenConnection            int    `yaml:"max_open_connection"`
	ConnectionMaxIdleTimeMinutes int    `yaml:"connection_max_idle_time_minutes"`
	ConnectionMaxLifeTimeMinutes int    `yaml:"connection_max_life_time_minutes"`
}

var (
	cfg  Config
	once sync.Once
)

var GetConfig = getConfigFunc

func getConfigFunc() Config {
	once.Do(func() {
		yamlFile, err := os.ReadFile("configs/cards.yaml")
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}

		// Parse the YAML data into a Config struct
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			log.Fatalf("Error parsing YAML: %v", err)
		}
	})

	return cfg
}
