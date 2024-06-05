package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const Usage string = `
	
	Usage:
		command-line: ./n3_fim baseline
			Get the filesystem baseline to monitor
		command-line: ./n3_fim monitor
			Start monitoring using the defined baseline
	`

const DefaultConfig string = `---

# Monitoring configs
monitoring:
  checking_frequency: 1 # Time between each check in seconds
  directories_to_monitor:
    - '.\examples_to_monitor'
  hashing_algorithm: "md5" # Supported algorithms: md5, sha256, crc32
  baseline_filename: "baseline.csv"
  log_file: "n3_fim.log"
`

// Global variable to hold the configuration
var config Config

// Config represents the overall configuration structure
type Config struct {
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

// MonitoringConfig represents the monitoring configuration section
type MonitoringConfig struct {
	CheckingFrequency    time.Duration `yaml:"checking_frequency"`
	DirectoriesToMonitor []string      `yaml:"directories_to_monitor"`
	HashingAlgorithm     string        `yaml:"hashing_algorithm"`
	BaselineFilename     string        `yaml:"baseline_filename"`
	LogFilename          string        `yaml:"log_file"`
}

func init() {
	// Read the YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the YAML data into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
