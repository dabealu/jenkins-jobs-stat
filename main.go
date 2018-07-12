package main

import (
	"flag"
)

func main() {
	confPath := flag.String("config", "config.yaml", "path to config file")
	example := flag.Bool("sample-config", false, "print example config and exit")

	flag.Parse()

	if *example {
		printConfExample()
	}

	conf := loadConf(*confPath)

	ch := make(chan Build)

	go getBuilds(conf, ch)

	saveBuilds(conf, ch)
}
