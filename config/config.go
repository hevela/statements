package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Interval       string
	StartAt        string
	FilesDir       string
	SendGridAPIKey string
	TemplateID     string
}

const configfile = "config/config.cfg"

var sampleConfig = `
# Indicates how often the process will work. E.g.: set to "730h" to run every 30 days
interval = "24h"
# Set startAt to empty to stop the worker. If you want an immediate start set it to "now", otherwise use a AM-PM hour format, e.g., 12:00PM
startAt = "12:00AM"
# Directory containing the statement CSV files, relative to the root of the project
filesDir = "statements"
# API key to be able to send emails
sendGridAPIKey = "add your sendgrid API key"
# the sendgrid dynamic template ID to use for mailing
templateID = "dynamic template ID"
`

func SampleConfig() string {
	return sampleConfig
}

func GetConfig() *Config {
	cfgFile, _ := filepath.Abs(configfile)
	conf := Config{}
	if _, err := toml.DecodeFile(cfgFile, &conf); err != nil {
		logrus.Error(err)
	}

	return &conf
}
