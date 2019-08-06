package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

const (
	AccessKeyId     = "aws_access_key_id"
	SecretAccessKey = "aws_secret_access_key"
	DefaultProfile  = "default"
)

func main() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		log.Fatalf("$HOME not set")
	}

	credentialsPath := flag.String("credentials", path.Join(homeDir, ".aws/credentials"), "path to the AWS credentials file")
	quiet := flag.Bool("q", false, "only display environment names")

	flag.Parse()

	var err error
	if len(flag.Args()) == 0 {
		err = printProfiles(*credentialsPath, *quiet)
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

func printProfiles(credentialsPath string, namesOnly bool) error {
	cfg, err := ini.Load(credentialsPath)
	if err != nil {
		return err
	}

	defaultKey, err := getValue(cfg, DefaultProfile, AccessKeyId)
	if err != nil {
		return err
	}

	for _, sec := range getProfiles(cfg) {
		key, err := getValue(cfg, sec, AccessKeyId)
		if err != nil {
			log.Printf("error reading profile %v: %v", sec, err)
			continue
		}

		if namesOnly {
			fmt.Println(sec)
		} else {
			if key == defaultKey {
				fmt.Printf("* %s\n", sec)
			} else {
				fmt.Printf("  %s\n", sec)
			}
		}
	}

	return nil
}

func setProfile(credentialsPath, profile string) error {
	cfg, err := ini.Load(credentialsPath)
	if err != nil {
		return err
	}

	err = saveDefaultProfile(cfg)
	if err != nil {
		return err
	}

	err = copyCredentials(cfg, profile, DefaultProfile)
	if err != nil {
		return err
	}

	err = cfg.SaveTo(credentialsPath)
	if err != nil {
		return err
	}

	fmt.Printf("switched to profile %s\n", profile)

	return nil
}

func saveDefaultProfile(cfg *ini.File) error {
	defaultKey, err := getValue(cfg, DefaultProfile, AccessKeyId)
	if err != nil {
		return err
	}

	// Check if the default credentials are stored in a named profile.
	found := false
	for _, sec := range getProfiles(cfg) {
		key, err := getValue(cfg, sec, AccessKeyId)
		if err != nil {
			log.Printf("error reading profile %v: %v", sec, err)
			continue
		}
		if defaultKey == key {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	// Create a new profile named "profileX".
	for i := 1; ; i++ {
		newProfile := fmt.Sprintf("profile%d", i)
		// Check if profile already exists.
		_, err := cfg.GetSection(newProfile)
		if err != nil {
			_, err := cfg.NewSection(newProfile)
			if err != nil {
				return err
			}

			err = copyCredentials(cfg, DefaultProfile, newProfile)
			if err != nil {
				return err
			}
			fmt.Printf("saved default credentials as profile %s\n", newProfile)
			break
		}
	}

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

func setValue(cfg *ini.File, profile, key, value string) error {
	sec, err := cfg.GetSection(profile)
	if err != nil {
		return err
	}

	sec.Key(key).SetValue(value)
	return nil
}

func getProfiles(cfg *ini.File) []string {
	allSections := cfg.SectionStrings()
	sections := []string{}
	for _, sec := range allSections {
		if sec != ini.DEFAULT_SECTION && sec != DefaultProfile {
			sections = append(sections, sec)
		}
	}
	return sections
}

func copyCredentials(cfg *ini.File, fromProfile, toProfile string) error {
	for _, key := range []string{AccessKeyId, SecretAccessKey} {
		value, err := getValue(cfg, fromProfile, key)
		if err != nil {
			return err
		}

		err = setValue(cfg, toProfile, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
