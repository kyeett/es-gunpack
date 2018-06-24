package main

import (
	"errors"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/kyeett/es-gunpack/pkg/unpacker"
	"github.com/urfave/cli"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/kyeett/es-gunpack/pkg/example-protofiles"
)

func main() {
	app := cli.NewApp()
	app.Name = "es-gunpack"
	app.Usage = `Golang client that monitors one or multiple Elasticsearch indices and unpacks binary data`
	app.HideHelp = true
	app.Version = "0.1"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Value: "localhost",
			Usage: "Url of the elasticsearch instance",
		},
		cli.StringFlag{
			Name:  "set-field-value",
			Usage: "Set the data field to value",
		},
		cli.BoolFlag{
			Name:  "reset-parsed",
			Usage: "Reset parsed-flag on document",
		},
		cli.BoolFlag{
			Name:  "set-parsed",
			Usage: "Set parsed-flag on document",
		},
		cli.BoolFlag{
			Name:  "ping",
			Usage: "Ping Elasticsearch instance. Tests authentication, url and port",
		},
	}
	app.Action = func(c *cli.Context) {

		url := "http://" + c.String("url") + ":9200"
		unpackerClient := unpacker.NewUnpacker(url, "logstash-2018.06.15")

		if c.Bool("ping") {
			unpackerClient.PingElasticsearch(url)
			os.Exit(0)
		}

		//Set tag parsed=false to all documeBnts
		if c.Bool("reset-parsed") {
			unpackerClient.SetParsedStatus(false)
			os.Exit(0)
		}

		//Set tag parsed=true to all documents
		if c.Bool("set-parsed") {
			unpackerClient.SetParsedStatus(true)
			os.Exit(0)
		}

		if c.String("set-field-value") != "" {
			unpackerClient.SetFieldStringValue("data", c.String("set-field-value"))
			os.Exit(0)
		}

		unpackerClient.ParseAndUpdate(printData)
	}
	app.Run(os.Args)
}

func printData(jsonMap map[string]interface{}) ([]byte, error) {

	if str, ok := jsonMap["data"].(string); ok {

		sDec, err := b64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, err
		}

		newTest := &example.Test{}
		err = proto.Unmarshal([]byte(sDec), newTest)
		if err != nil {
			return nil, err
		}

		j, err := json.Marshal(newTest)
		if err != nil {
			return nil, err
		}
		return j, nil

	} else {
		return nil, errors.New("Field 'data' is not string and cannot be decoded")
	}

}
