package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type checkHTTP struct {
	HTTP []cHTTP `json:"http"`
}

type cHTTP struct {
	Hostname string `json:"hostname"`
	Enable   bool   `json:"enable"`
	Timeout  int    `json:"timeout"`
	Port     int    `json:"port"`
	SSL      bool   `json:"bool"`
	Data     int    `json:"data"`
	Status   int    `json:"status"`
}

func (c *cHTTP) check() error {
	client := http.Client{
		Timeout: time.Second * time.Duration(c.Timeout),
	}

	url := "http://"
	if c.SSL {
		url = "https://"
	}

	sPort := strconv.Itoa(c.Port)
	resp, err := client.Get(url + c.Hostname + ":" + sPort)
	if err != nil {
		return err
	}

	// add outout to struct
	c.Data = resp.StatusCode

	// return error
	return err
}

func (c *cHTTP) status() (bool, string) {
	if c.Data != c.Status {
		return false, convertData(c.Data)
	}
	return true, ""
}

func setHTTP() (checkHTTP, error) {
	f, err := getConfig("check_http")
	if err != nil {
		return checkHTTP{}, err
	}
	if f == "" {
		errorf("Cannot find config file for check:", f)
	}

	httpST := checkHTTP{}
	fileByte, err := readConf(f)
	if err != nil {
		return checkHTTP{}, err
	}

	err = json.Unmarshal(fileByte, &httpST)
	if err != nil {
		return checkHTTP{}, err
	}

	return httpST, err
}
