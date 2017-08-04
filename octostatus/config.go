package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

/*Config describes the config file format for octostatus

Top level fields
Printer and Output map keys currently in use is "default"
*/
type Config struct {
	Name    string                   `yaml:"name"`
	Version int16                    `yaml:"version"`
	Printer map[string]PrinterConfig `yaml:"printer"`
	Output  map[string]OutputConfig  `yaml:"output"`
}

/*PrinterConfig Octoprint URL and Credentials

Mandatory for usage
*/
type PrinterConfig struct {
	BaseURL string `yaml:"baseUrl"`
	APIKey  string `yaml:"apiKey"`
}

/*OutputConfig Template(s) for output generation

Mandatory for usage
*/
type OutputConfig struct {
	Template string `yaml:"template"`
}

func readConfig(filename string) (Config, error) {
	filepath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if config.Version != 1 {
		myerr := fmt.Errorf("Unknown configuration version: %d", config.Version)
		return config, myerr
	}
	return config, nil
}
