package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic"
	"github.com/urfave/cli"
)

//
// Unpacker is a structure that holds data necessary to monitor and unpack data in Elasticsearch
//
type Unpacker struct {
	client   *elastic.Client // elasticsearch client
	indicies []string
}

func (u Unpacker) setParsedStatus(status bool) {
	fmt.Println("All tags reset\n")

	termQuery := elastic.NewMatchAllQuery()
	ctx := context.Background()
	_, err := u.client.UpdateByQuery(u.indicies[0]).
		Query(termQuery).
		Script(elastic.NewScript("ctx._source.parsed = params.status").Param("status", status)).
		//		Script(elastic.NewScript("ctx._source.parsed = true").Params(map[string]interface{}{"tag": "blue"}).Lang("painless")).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

type document struct {
	to    string
	from  string
	event string
}

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

		unpacker := Unpacker{client: client, indicies: []string{"logstash-2018.06.15"}}

		//Set tag parsed=false to all documents
		if c.Bool("reset-parsed") {
			unpacker.setParsedStatus(false)
			os.Exit(0)
		}

		//Set tag parsed=true to all documents
		if c.Bool("set-parsed") {
			unpacker.setParsedStatus(true)
			os.Exit(0)
		}

		boolTermQuery := elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("parsed", true))
		/*
		   {
		     "query": {
		       "bool" : {
		         "must_not" : {
		           "term" : {
		             "parsed" : true
		           }
		         }
		       }
		     }
		   }
		*/

		var updateIndex string = "logstash-2018.06.15"
		result, err := client.Search().
			Index(updateIndex).
			From(0).
			Size(9000). //TODO: needs rewrite this using scrolling, as this implementation may loose entries if there's more than 9K entries per sleep period
			Query(boolTermQuery).
			Do(context.Background())

		if err != nil {
			// Handle error
			panic(err)
		}

		// result is of type result and returns hits, suggestions,
		// and all kinds of other information from Elasticsearch.
		fmt.Printf("Query took %d milliseconds\n", result.TookInMillis)

		// Here's how you iterate through results with full control over each step.
		if result.Hits.TotalHits > 0 {
			fmt.Printf("\nFound a total of %d unparsed signals\n", result.Hits.TotalHits)

			jsonMap := make(map[string]interface{})

			// Iterate through results
			for _, hit := range result.Hits.Hits {
				// hit.Index contains the name of the index

				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				err := json.Unmarshal(*hit.Source, &jsonMap)
				if err != nil {
					// Deserialization failed
					fmt.Printf("Deserialization failed %v\n", err)
				}

				fmt.Printf("Update document with ID %v\n", hit.Id)
				// update, _ := client.Update().Index(updateIndex).Type("doc").Id(hit.Id).
				// 	Script(elastic.NewScriptInline("ctx._source.parsed = false").Lang("painless")).
				// 	Do(ctx)
			}
		} else {
			fmt.Print("\nFound no unparsed signals found\n")
		}

	}
	app.Run(os.Args)
}
