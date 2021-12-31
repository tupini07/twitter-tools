package app_config

import (
	"io/ioutil"
	"log"
	"os"
)

func writeDefaultConfig() {
	defaultConfig := `
auth:
  api_key: put your api_key here
  api_secret_key: put your api_secret_key here
  
  access_token: put your access_token here
  access_token_secret: put your access_token_secret here
  
log_level: INFO

flow:
repeat: true
steps:
    - follow_all_followers:
        max_to_follow: 200

    - follow_followers_of_others:
        max_to_follow: 100
        random: false
        others:
        - list of twitter handlers

    - unfollow_bad_friends:
        max_to_unfollow: 400

    - wait:
        seconds: 27
        minutes: 30
        hours: 0
      
`

	if _, err := os.Stat("config.yml"); err == nil {
		log.Fatal("Trying to write default config file but it already exists!")
	}

	ioutil.WriteFile("config.yml", []byte(defaultConfig), 0770)
}
