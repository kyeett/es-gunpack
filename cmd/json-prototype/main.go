package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kyeett/es-gunpack/pkg/unpacker"
	"github.com/olivere/elastic"
)

// Simple cli that adds a protobuf signal to the 'data' field of all entries in logstash-all
func main() {

	// Create client interfacing elasticsearch
	url := "http://localhost:9200"
	unpackerClient := unpacker.NewUnpacker(url, "logstash-2018.06.15")

	boolTermQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("parsed", true))
	result, err := unpackerClient.Client.Search().
		Index(unpackerClient.Index).
		From(0).
		Size(10000).
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
		for _, hit := range result.Hits.Hits {
			jsonMap := make(map[string]interface{})

			// Deserialize hit.Source
			err := json.Unmarshal(*hit.Source, &jsonMap)
			if err != nil {
				// Deserialization failed
				fmt.Printf("Deserialization failed %v\n", err)
			}
			fmt.Printf("%v\n", jsonMap["unpacked"])

			if str, ok := jsonMap["unpacked"].(string); ok {
				sDec, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%v\n", string(sDec))
			}
		}
	}
}
