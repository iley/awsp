package main

import (
	"flag"
	"fmt"
	"os"
	// "gopkg.in/ini.v1"
)

var credentialsPath = flag.String("credentials", "~/.aws/credentials", "path to the AWS credentials file")
var configPath = flag.String("config", "~/.aws/config", "path to the AWS config file")
var proflie = flag.String("set", "", "profile to set as default")

func main() {
	flag.Parse()

	var err error
	if *proflie != "" {
		err = setProfile()
	} else {
		err = getProfiles()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getProfiles() error {
	fmt.Println("getProfiles()")
	return nil
}

func setProfile() error {
	fmt.Println("setProfile()")
	return nil
}
