package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const langConfigFname = ".available_langs.json"

type LangConfig []struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Get config path
func getFullConfigPath(cfgFpath string) (fPath string, err error) {
	fPath, err = filepath.Abs(".")
	if err != nil {
		log.Printf("Error getting absolute path to current dir: %s", err)
		return
	}

	return filepath.Join(fPath, cfgFpath), err
}

// Read locale-related config file
func ReadLangs() (cfg LangConfig, err error) {
	fmt.Println("Ran")
	path, err := getFullConfigPath(langConfigFname)
	if err != nil {
		log.Printf("Error getting full config path: %s", err)
		return
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file at %s: %s", path, err)
		return
	}

	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		log.Printf("Error unmarshalling file content: %s", err)
	}

	return cfg, err
}
