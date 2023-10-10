package main

import (
	"Client/Structs"
	b64 "encoding/base64"
	"encoding/json"
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

	client := Structs.NewClient(config)
	sendVersion(err, client)
	sendDecode(err, client)
	sendHardOp(err, client)
}

func sendHardOp(err error, client *Structs.Client) {
	data, err := client.SendLongRequest("hard-op")
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(data)
}

func sendVersion(err error, client *Structs.Client) {
	data, err := client.SendGetRequest("version")
	if err != nil {
		log.Error(err)
		return
	}

	var response Structs.Response
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(response.Message)
}

func sendDecode(err error, client *Structs.Client) {
	encoded := b64.StdEncoding.EncodeToString([]byte("Hello world! :)"))
	jsonRequest, err := json.Marshal(Structs.Request{Message: encoded})
	if err != nil {
		log.Error(err)
		return
	}

	data, err := client.SendPostRequest("decode", jsonRequest)
	if err != nil {
		log.Error(err)
		return
	}

	var response Structs.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(response.Message)
}
