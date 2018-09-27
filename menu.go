package main

import "flag"

var menu argstruct

type argstruct struct {
	configfile string
}

func flagOptions() argstruct {
	configfile := flag.String("config", "", "config file")
	flag.Parse()

	return argstruct{*configfile}
}
