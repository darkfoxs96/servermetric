package telegram

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	config *Config
)

type Config struct {
	Subs  []int64 // List of subscribers to notifications
	mutex sync.Mutex
}

var configFileName = "tconfig.cfg"

func start(pathConf string) {
	configFileName = pathConf
	config = loadConfig()
}

func loadConfig() *Config {
	var conf Config
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Println("can't load data:", err)
		return &Config{}
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Println("can't unmarshal data:", err)
		return &Config{}
	}
	return &conf
}

func (conf *Config) save() {
	conf.mutex.Lock()
	defer conf.mutex.Unlock()

	data, err := json.Marshal(conf)
	if err != nil {
		log.Println("can't marshal config:", err)
		return
	}
	err = ioutil.WriteFile(configFileName, data, 0666)
	if err != nil {
		log.Println("can't save config:", err)
		return
	}
}
