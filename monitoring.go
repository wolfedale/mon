package main

import "fmt"

type monitoring interface {
	check() error
	status() (bool, string)
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

func convertData(i interface{}) string {
	return fmt.Sprint(i)
}
