package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

func main() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		log.Fatalf("$HOME not set")
	}

	credentialsPath := flag.String("credentials", path.Join(homeDir, ".aws/credentials"), "path to the AWS credentials file")
	_ = flag.String("config", path.Join(homeDir, ".aws/config"), "path to the AWS config file")
	proflie := flag.String("set", "", "profile to set as default")

	flag.Parse()

	var err error
	if *proflie != "" {
		err = setProfile()
	} else {
		err = getProfiles(*credentialsPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getProfiles(credentialsPath string) error {
	cfg, err := ini.Load(credentialsPath)
	if err != nil {
		return err
	}

	defaultSec, err := cfg.GetSection("default") // not to be confused with DEFAULT which is top-level INI section
	if err != nil {
		return err
	}
	if !defaultSec.HasKey("aws_access_key_id") {
		return fmt.Errorf("aws_access_key_id not found in \"default\" section of %s", credentialsPath)
	}

	defaultKey := defaultSec.Key("aws_access_key_id").Value()

	for _, sec := range cfg.Sections() {
		if sec.Name() == "DEFAULT" || sec.Name() == "default" {
			continue
		}
		key := sec.Key("aws_access_key_id").Value()
		if key != "" && key == defaultKey {
			fmt.Printf("* %s\n", sec.Name())
		} else {
			fmt.Printf("  %s\n", sec.Name())
		}
	}

	return nil
}

func setProfile() error {
	fmt.Println("setProfile()")
	return nil
}
