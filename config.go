package main

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Config []conf `yaml:"config"`
}

type conf struct {
	Name string `yaml:"name"`
	Conf string `yaml:"conf"`
}

func getConfig(name string) (string, error) {
	var c config
	source, err := ioutil.ReadFile(menu.configfile)
	if err != nil {
		errorf("Cannot read file: %v", err)
	}
	err = yaml.Unmarshal(source, &c)
	if err != nil {
		errorf("Cannot Unmarshal yaml file: %v", err)
	}
	var n string
	for _, j := range c.Config {
		if j.Name == name {
			return j.Conf, nil
		}
	}
	return n, nil
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
