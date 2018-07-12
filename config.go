package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// common config parameters
type Common struct {
	Interval int `yaml:"interval"`
}

// jenkins access credentials and other parameters
type Jenkins struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Url      string `yaml:"url"`
}

// access credentials to postgres database
type Datastore struct {
	Url      string `yaml:"url"`
	Index    string `yaml:"index"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// full configuration represented by config file
type config struct {
	Common    `yaml:"common"`
	Jenkins   `yaml:"jenkins"`
	Datastore `yaml:"datastore"`
}

// print config sample
func printConfExample() {
	c := config{
		Common{
			Interval: 600,
		},
		Jenkins{
			User:     "user",
			Password: "password",
			Url:      "http://localhost:8080",
		},
		Datastore{
			Url:      "http://localhost:9200",
			Index:    "jenkins",
			User:     "user",
			Password: "password",
		},
	}

	b, err := yaml.Marshal(c)
	logFatal("Print config example", err)

	fmt.Println(string(b))

	os.Exit(0)
}

// load configuration file
func loadConf(path string) *config {
	conf := &config{}

	file, err := os.Open(path)
	defer file.Close()
	logFatal("Load conf: open file", err)

	b, err := ioutil.ReadAll(file)
	logFatal("Load conf: read file", err)

	err = yaml.Unmarshal(b, conf)
	logFatal("Load conf: json", err)

	log.Printf("Config loaded: %s", path)

	return conf
}
