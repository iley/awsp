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

	flag.Parse()

	var err error
	if len(flag.Args()) == 0 {
		err = getProfiles(*credentialsPath)
	} else {
		// TODO: Update ~/.aws/config as well.
		profile := flag.Args()[0]
		err = setProfile(*credentialsPath, profile)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getProfiles(credentialsPath string) error {
	cfg, err := ini.Load(credentialsPath)
	if err != nil {
		return err
	}

	defaultKey, err := getValue(cfg, "default", "aws_access_key_id")
	if err != nil {
		return err
	}

	for _, sec := range cfg.SectionStrings() {
		if sec == "DEFAULT" || sec == "default" {
			continue
		}

		key, err := getValue(cfg, sec, "aws_access_key_id")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		if key == defaultKey {
			fmt.Printf("* %s\n", sec)
		} else {
			fmt.Printf("  %s\n", sec)
		}
	}

	return nil
}

func setProfile(credentialsPath, profile string) error {
	cfg, err := ini.Load(credentialsPath)
	if err != nil {
		return err
	}

	for _, key := range []string{"aws_access_key_id", "aws_secret_access_key"} {
		value, err := getValue(cfg, profile, key)
		if err != nil {
			return err
		}

		err = setValue(cfg, "default", key, value)
		if err != nil {
			return err
		}
	}

	err = cfg.SaveTo(credentialsPath)
	if err != nil {
		return err
	}

	fmt.Printf("switched to profile %s\n", profile)

	return nil
}

func getValue(cfg *ini.File, section, key string) (string, error) {
	sec, err := cfg.GetSection(section)
	if err != nil {
		return "", err
	}

	if !sec.HasKey(key) {
		return "", fmt.Errorf("%s not found in section \"%s\"", key, section)
	}

	value := sec.Key(key).Value()
	if value == "" {
		return "", fmt.Errorf("%s empty in section \"%s\"", key, section)
	}

	return value, nil
}

func setValue(cfg *ini.File, section, key, value string) error {
	sec, err := cfg.GetSection("default")
	if err != nil {
		return err
	}

	sec.Key(key).SetValue(value)
	return nil
}
