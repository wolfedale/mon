package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	menu = flagOptions()
	if menu.configfile == "" {
		flag.PrintDefaults()
		errorf("Config not specified\n")
	}

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

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}
