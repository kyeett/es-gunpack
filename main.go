package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/olivere/elastic.v5"
)

//
// Unpacker is a structure that holds data necessary to monitor and unpack data in Elasticsearch
//
type Unpacker struct {
	client   *elastic.Client // elasticsearch client
	indicies []string
}

func main() {
	app := cli.NewApp()
	app.Name = "es-gunpack"
	app.Usage = `Golang client that monitors one or multiple Elasticsearch indices and unpacks binary data`
	app.HideHelp = true
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}
	app.Action = func(c *cli.Context) {
		fmt.Println(c.String("lang"))
		ctx := context.Background()

		defaultOptions := []elastic.ClientOptionFunc{
			elastic.SetURL("http://localhost:9200"),
			elastic.SetSniff(false),
			elastic.SetBasicAuth("elastic", "changeme"),
			elastic.SetHealthcheckTimeoutStartup(10 * time.Second),
			elastic.SetHealthcheckTimeout(2 * time.Second),
		}

		client, err := elastic.NewClient(defaultOptions...)

		if err != nil {
			// Handle error
			log.Println("first")
			log.Fatal(err)
		}

		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
		if err != nil {
			// Handle error
			log.Fatal(err)
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		termQuery := elastic.NewMatchAllQuery()
		result, err := client.Search().
			Index("logstash-*").
			From(0).
			Size(9000). //TODO: needs rewrite this using scrolling, as this implementation may loose entries if there's more than 9K entries per sleep period
			Query(termQuery).
			Do(context.Background())

		if err != nil {
			// Handle error
			panic(err)
		}

		// result is of type result and returns hits, suggestions,
		// and all kinds of other information from Elasticsearch.
		fmt.Printf("Query took %d milliseconds\n", result.TookInMillis)
		fmt.Printf("%+v\n", result)

		// Here's how you iterate through results with full control over each step.
		if result.Hits.TotalHits > 0 {
			fmt.Printf("Found a total of %d tweets\n", result.Hits.TotalHits)

			/*
			   // Iterate through results
			   for _, hit := range result.Hits.Hits {
			      // hit.Index contains the name of the index

			      // Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			      var t Tweet
			        err := json.Unmarshal(*hit.Source, &t)
			        if err != nil {
			           // Deserialization failed
			        }

			        // Work with tweet
			        fmt.Printf("Tweet by %s: %s\n", t.To, t.From)
			   }
			*/

		} else {
			// No hits
			fmt.Print("Found no tweets\n")
		}

	}
	app.Run(os.Args)
}
