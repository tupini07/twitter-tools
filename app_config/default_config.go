package app_config

import (
	_ "embed"
	"io/ioutil"
	"os"

	"github.com/tupini07/twitter-tools/print_utils"
)

//go:embed config.example.yml
var defaultConfig string

func writeDefaultConfig() {

	if _, err := os.Stat("config.yml"); err == nil {
		print_utils.Fatal("Trying to write default config file but it already exists!")
	}

	ioutil.WriteFile("config.yml", []byte(defaultConfig), 0770)
}
