package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

type monitoring interface {
	check() error
	status() (bool, string)
}

func main() {
	icmp, err := setICMP()
	if err != nil {
		fmt.Printf("Error when setting ICMP struct: %s\n", err)
	}
	http, err := setHTTP()
	if err != nil {
		fmt.Printf("Error when setting HTTP struct: %s\n", err)
	}

	// run icmp checks on every host
	for _, h := range icmp.ICMP {
		err := runChecks(&h)
		if err != nil {
			fmt.Printf("Error when executing check(): %s\n", err)
		}
	}

	// run http checks on every host
	for _, h := range http.HTTP {
		err := runChecks(&h)
		if err != nil {
			fmt.Printf("Error when executing check(): %s\n", err)
		}
	}

	// run N checks
}

func runChecks(m monitoring) error {
	// run checks
	err := m.check()
	if err != nil {
		return err
	}

	// get status
	status, data := m.status()
	fmt.Println(status, data)
	return err
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

func (c *cICMP) status() (bool, string) {
	if c.Data > c.Max {
		return false, convertData(c.Data)
	}
	return true, ""
}

func (c *cHTTP) status() (bool, string) {
	if c.Data != c.Status {
		return false, convertData(c.Data)
	}
	return true, ""
}

func convertData(i interface{}) string {
	return fmt.Sprint(i)
}

func setICMP() (checkICMP, error) {
	icmpST := checkICMP{}

	fileByte, err := readConf("./tests/icmp.yaml")
	if err != nil {
		return checkICMP{}, err
	}

	err = json.Unmarshal(fileByte, &icmpST)
	if err != nil {
		return checkICMP{}, err
	}

	return icmpST, nil
}

func setHTTP() (checkHTTP, error) {
	httpST := checkHTTP{}

	fileByte, err := readConf("./tests/http.yaml")
	if err != nil {
		return checkHTTP{}, err
	}

	err = json.Unmarshal(fileByte, &httpST)
	if err != nil {
		return checkHTTP{}, err
	}

	return httpST, err
}

func readConf(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	byteFile, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}

	return byteFile, err
}
