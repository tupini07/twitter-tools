package app_config

import (
	_ "embed"
	"io/ioutil"
	"log"
	"os"
)

//go:embed config.example.yml
var defaultConfig string

func writeDefaultConfig() {

	if _, err := os.Stat("config.yml"); err == nil {
		log.Fatal("Trying to write default config file but it already exists!")
	}

	ioutil.WriteFile("config.yml", []byte(defaultConfig), 0770)
}
