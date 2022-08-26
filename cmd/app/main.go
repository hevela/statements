package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/hevela/statements/config"
	"github.com/hevela/statements/internal/app"
)

var cfg *config.Config

func main() {
	// call some global functions to get things set up
	if !getConfig() {
		logrus.Error("configuration error, check that all required values are set")
		return
	}
	sampleconfig := flag.Bool("sampleconfig", false, "Outputs a sample config to standard out")
	flag.Parse()

	if *sampleconfig {
		fmt.Println(config.SampleConfig())
		return
	}
	app.Run(cfg)
}

func getConfig() bool {
	cfg = config.GetConfig()
	if cfg == nil {
		logrus.Error("Unable to get config")
		return false
	}
	return true
}
