package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/darkfoxs96/servermetric/alert"
	"github.com/darkfoxs96/servermetric/db"
	"github.com/darkfoxs96/servermetric/pusher"
	_ "github.com/darkfoxs96/servermetric/pusher/pushers/telegram"
	"github.com/darkfoxs96/servermetric/web"
)

type GlobalConfig struct {
	AlertConfig *alert.AlertConfig
	Pushers     map[string]map[string]interface{}
	Database    *db.DBConfig
	Metrics     map[string]map[string]string
	Server      *web.ServerConfig
}

func main() {
	config := &GlobalConfig{}
	configFilePath := os.Getenv("SERVERMETRICCONF")

	if configFilePath == "" {
		configFilePath = "/servermetric.config.yml"
	}

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		panic(err)
	}

	// Database
	db.Run(config.Database)
	if err = db.AddMetrics(config.Metrics); err != nil {
		panic(err)
	}

	// Pushers
	if err = pusher.Run(config.Pushers); err != nil {
		panic(err)
	}

	// Alert
	alert.Run(config.AlertConfig)

	// API
	go web.Run(config.Server)
}
