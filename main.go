package main

import (
	"flag"
)

var conf *config

func init() {
	confPath := flag.String("config", "config.yaml", "path to config file")
	example := flag.Bool("sample-config", false, "print example config and exit")

	flag.Parse()

	if *example {
		printConfExample()
	}

	conf = loadConf(*confPath)

	datClient = createHttpClient(conf.Datastore.Verify)
	jClient = createHttpClient(conf.Jenkins.Verify)
}

func main() {
	ch := make(chan Build)

	go getBuilds(conf, ch)

	saveBuilds(conf, ch)
}
