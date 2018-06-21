package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kyeett/es-gunpack/pkg/unpacker"
	"github.com/olivere/elastic"
	"github.com/urfave/cli"
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

		ctx := context.Background()
		defaultOptions := []elastic.ClientOptionFunc{
			elastic.SetURL(url),
			elastic.SetSniff(false),
			elastic.SetBasicAuth("elastic", "changeme"),
			elastic.SetHealthcheckTimeoutStartup(10 * time.Second),
			elastic.SetHealthcheckTimeout(2 * time.Second),
		}

		client, err := elastic.NewClient(defaultOptions...)

		if err != nil {
			// Handle error
			log.Fatal(err)
		}

		if c.Bool("ping") {
			info, code, err := client.Ping(url).Do(ctx)
			if err != nil {
				// Handle error
				log.Fatal(err)
			}
			fmt.Printf("\nElasticsearch returned OK with code %d and version %s\n\n", code, info.Version.Number)
			os.Exit(0)
		}

		unpackerClient := unpacker.Unpacker{client, "logstash-2018.06.15"}

		//Set tag parsed=false to all documents
		if c.Bool("reset-parsed") {
			unpackerClient.SetParsedStatus(false)
			os.Exit(0)
		}

		//Set tag parsed=true to all documents
		if c.Bool("set-parsed") {
			unpackerClient.SetParsedStatus(true)
			os.Exit(0)
		}

		unpackerClient.ParseAndUpdate()
	}
	app.Run(os.Args)
}
