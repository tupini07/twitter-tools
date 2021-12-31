package main

import (
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/cmd"
	"github.com/tupini07/twitter-tools/database"
)

func initLogger() {
	conf := app_config.GetConfig()

	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		PadLevelText:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})

	switch conf.LogLevel {
	case "DEBUG":
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	case "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.WithFields(log.Fields{
			"log_level": conf.LogLevel,
		}).Warn("Setting log level to INFO. Unknown provided level")
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())

	initLogger()
	database.InitDb()

	cmd.RunCli()
}
