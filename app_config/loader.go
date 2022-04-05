package app_config

import (
	"io/ioutil"

	"github.com/tupini07/twitter-tools/print_utils"
	"gopkg.in/yaml.v2"
)

type FlowStep struct {
	Random *struct {
		Options []FlowStep `yaml:"options"`
	} `yaml:"random,omitempty"`
	FollowAllFollowers *struct {
		MaxToFollow int `yaml:"max_to_follow"`
	} `yaml:"follow_all_followers,omitempty"`
	FollowFollowersOfOthers *struct {
		MaxToFollow      int      `yaml:"max_to_follow"`
		MaxSourcesToPick int      `yaml:"max_sources_to_pick"`
		Others           []string `yaml:"others"`
	} `yaml:"follow_followers_of_others,omitempty"`
	UnfollowBadFriends *struct {
		MaxToUnfollow int `yaml:"max_to_unfollow"`
	} `yaml:"unfollow_bad_friends,omitempty"`
	Wait *struct {
		Seconds int64 `yaml:"seconds"`
		Minutes int64 `yaml:"minutes"`
		Hours   int64 `yaml:"hours"`
	} `yaml:"wait,omitempty"`
	WaitUntilDay *struct {
		Relative string `yaml:"relative"`
	} `yaml:"wait_until_day,omitempty"`
}

type Flow struct {
	Repeat            bool       `yaml:"repeat"`
	MaxTotalFollowing int        `yaml:"max_total_following"`
	Steps             []FlowStep `yaml:"steps"`
}

type AppConfig struct {
	Auth struct {
		APIKey            string `yaml:"api_key"`
		APISecretKey      string `yaml:"api_secret_key"`
		AccessToken       string `yaml:"access_token"`
		AccessTokenSecret string `yaml:"access_token_secret"`
	} `yaml:"auth"`
	LogLevel string `yaml:"log_level"`
	Flow     *Flow  `yaml:"flow"`
}

func readConfigFile() string {
	content, err := ioutil.ReadFile("config.yml")
	if err != nil {
		writeDefaultConfig()
		print_utils.Fatal("No config file was found, so a default one has been writter to the current directory. Please modify values as desired and try again. Exiting..")
	}

	return string(content)
}

var configInstance *AppConfig

func GetConfig() *AppConfig {
	if configInstance == nil {
		data := readConfigFile()

		c := AppConfig{}

		err := yaml.Unmarshal([]byte(data), &c)
		if err != nil {
			print_utils.Fatalf("error parsing config file: %v", err)
		}

		configInstance = &c
	}

	return configInstance
}
