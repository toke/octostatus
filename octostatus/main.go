package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	home := os.Getenv("HOME")
	cfg, err := readConfig(home + "/.config/octoclient/config.yml")
	if err != nil {
		panic(err)
	}

	printerID := "default"
	templateID := "default"

	var URL = cfg.Printer[printerID].BaseURL + "/api/job"
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
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	req.Header.Set("X-Api-Key", APIKEY)
	response, err := netClient.Do(req)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
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
		err = tmpl.Execute(os.Stdout, t)
		fmt.Println("")
	}

}
