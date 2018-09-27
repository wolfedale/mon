package main

import (
	"encoding/json"
	"time"

	ping "github.com/sparrc/go-ping"
)

type checkICMP struct {
	ICMP []cICMP `json:"icmp"`
}

type cICMP struct {
	Hostname string  `json:"hostname"`
	Enable   bool    `json:"enable"`
	Timeout  int     `json:"timeout"`
	Data     float64 `json:"data"`
	Max      float64 `json:"max"`
}

func (c *cICMP) check() error {
	pinger, err := ping.NewPinger(c.Hostname)
	if err != nil {
		return err
	}

	pinger.Count = c.Timeout
	pinger.Timeout = time.Duration(c.Timeout) * time.Second
	pinger.Run()

	// get stats
	stats := pinger.Statistics()

	// add output to struct
	c.Data = stats.PacketLoss

	// return error
	return err
}

func (c *cICMP) status() (bool, string) {
	if c.Data > c.Max {
		return false, convertData(c.Data)
	}
	return true, ""
}

func setICMP() (checkICMP, error) {
	f, err := getConfig("check_icmp")
	if err != nil {
		return checkICMP{}, err
	}
	if f == "" {
		errorf("Cannot find config file for check:", f)
	}

	icmpST := checkICMP{}
	fileByte, err := readConf(f)
	if err != nil {
		return checkICMP{}, err
	}
	err = json.Unmarshal(fileByte, &icmpST)
	if err != nil {
		return checkICMP{}, err
	}

	return icmpST, nil
}
