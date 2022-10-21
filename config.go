package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Token    string `koanf:"token"`
	Zones    []Zone `koanf:"zones"`
	Interval int    `koanf:"interval"`
}

type Zone struct {
	Name  string   `koanf:"zone"`
	TTL   int      `koanf:"ttl"`
	Hosts []string `koanf:"hosts"`
}

func ReadConfig(configFile string) (conf *Config, err error) {
	f := file.Provider(configFile)
	k := koanf.New("/")
	// load default values
	k.Load(confmap.Provider(map[string]interface{}{
		"cfd/interval": 300,
	}, "/"), nil)
	// load YAML config
	if err = k.Load(f, yaml.Parser()); err != nil {
		return
	}
	// unmarshal the data
	err = k.Unmarshal("cfd", &conf)
	if err != nil {
		return
	}
	// watch the file and get a callback on change
	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Fatalf("file watch error: %v", err)
		}
		// throw away the old config and load a fresh copy
		log.Println("Config changed. Reloading...")
		k = koanf.New("/")
		if err = k.Load(f, yaml.Parser()); err != nil {
			log.Fatalf("Unable to load config file: %v", err)
			return
		}
		if err = k.Unmarshal("cfd", &conf); err != nil {
			log.Fatalf("Inavlid change to config file: %v", err)
			return
		}
	})
	log.Printf("loaded config from %s", configFile)
	return
}
