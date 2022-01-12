package main

// Handles writing and reading world data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/yaml.v3"
)

var num_platforms = 0

// Writes the given platform to the file at path
func WritePlatformData(p Platform, path string) {
	platform_name, err := uuid.NewV4()
	platform := map[string]Platform{platform_name.String(): p}
	data, err1 := yaml.Marshal(&platform)
	if err1 != nil {
		log.Fatal(err)
	}
	f, err2 := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}
	if _, err := f.Write(data); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// Reads every line from the given yaml file and returns an aray of platforms.
func ReadPlatformData(path string) []Platform {
	data_file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	result_data := make(map[string]Platform)
	err1 := yaml.Unmarshal(data_file, result_data)
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Printf("\n%v\n", result_data)
	var result []Platform
	for _, value := range result_data {
		result = append(result, value)
	}

	return result
}
