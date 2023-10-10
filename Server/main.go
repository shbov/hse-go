package main

import (
	"Server/Structs"
	"flag"
	log "github.com/sirupsen/logrus"
)

func main() {
	path := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := Structs.ParseConfig(*path)
	if err != nil {
		log.Error(err)
		return
	}

	server := Structs.NewServer(config)
	server.Start()
}
