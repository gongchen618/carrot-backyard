package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	C *Config
)

type Config struct {
	App     app     `yaml:"app"`
	MongoDB mongodb `yaml:"mongodb"`
	QQBot   qqbot   `yaml:"qqbot"`
}

type mongodb struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type app struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}

type qqbot struct {
	Host string `yaml:"host"`
}

func init() {
	configFile := "config/default.yml"

	// 如果设置了
	if v, ok := os.LookupEnv("CONFIG"); ok {
		configFile = "config/" + v + ".yml"
	}

	configFilePath := ""
	p, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		return
	}

	for {
		configFilePath = path.Join(p, configFile)
		if _, err := os.Stat(configFilePath); err == nil {
			break
		}
		if p == "/" {
			log.Panic("config not found")
		}
		p = path.Join(p, "..")
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("Read config error!")
		log.Panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	C = config

	log.Println("Config " + configFile + " loaded.")
}
