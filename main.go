package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type JobInfo struct {
	JobDetail JobDetail   `json:"job"`
	Progress  JobProgress `json:"progress"`
	State     string      `json:"state"`
}

type JobDetail struct {
	AveragePrintTime   float64                `json:"averagePrintTime"`
	EstimatedPrintTime float64                `json:"estimatedPrintTime"`
	LastPrintTime      float64                `json:"lastPrintTime"`
	File               JobFile                `json:"file"`
	Tool               map[string]JobFilament `json:"filament"`
}

type JobFile struct {
	Date   int    `json:"date"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Path   string `json:"path"`
	Size   int32  `json:"size"`
}

type JobFilament struct {
	Length float64 `json:"length"`
	Volume float64 `json:"volume"`
}

type JobProgress struct {
	ETA                 string  `json:"ETA,omitempty"`
	Completion          float64 `json:"completion,omitempty"`
	Filepos             int32   `json:"filepos"`
	PrintTime           int32   `json:"printTime"`
	PrintTimeLeft       int32   `json:"printTimeLeft"`
	PrintTimeLeftOrigin string  `json:"printTimeLeftOrigin,omitempty"`
	PrintTimeLeftString string  `json:"printTimeLeftString"`
}

type Config struct {
	Name    string                   `yaml:"name"`
	Version int16                    `yaml:"version"`
	Printer map[string]PrinterConfig `yaml:"printer"`
	Output  map[string]OutputConfig  `yaml:"output"`
}

type PrinterConfig struct {
	BaseURL string `yaml:"baseUrl"`
	URL     string `yaml:"url"`
	APIKey  string `yaml:"apiKey"`
}

type OutputConfig struct {
	Template string `yaml:"template"`
}

func readConfig(filename string) Config {
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
	return (config)
}

func main() {
	home := os.Getenv("HOME")
	cfg := readConfig(home + "/.config/octoclient/config.yml")

	printerID := "default"
	templateID := "default"

	var URL = cfg.Printer[printerID].URL
	var APIKEY = cfg.Printer[printerID].APIKey
	if URL == "" || APIKEY == "" {
		fmt.Printf("Error: Missing URL or APIKEY for printer \"%s\"\n", printerID)
		os.Exit(1)
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	req.Header.Set("X-Api-Key", APIKEY)
	response, err := netClient.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		//fmt.Printf("%s\n", contents)

		var t JobInfo
		json.Unmarshal(contents, &t)

		dur := time.Duration(t.Progress.PrintTimeLeft) * time.Second
		if t.Progress.PrintTimeLeftString == "" {
			t.Progress.PrintTimeLeftString = dur.String()
		}

		tmpl, err := template.New("output").Parse(cfg.Output[templateID].Template)
		//tmpl, err := template.New("test").Parse("{{.State}} {{.JobDetail.File.Name}} {{printf \"%.01f\" .Progress.Completion }}% ({{ .Progress.PrintTimeLeftString }})\n")
		err = tmpl.Execute(os.Stdout, t)
		fmt.Println("")
	}

}
