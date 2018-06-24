package unpacker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/olivere/elastic"
)

//
// unpacker is a structure that holds data necessary to monitor and unpack data in Elasticsearch
//
type unpacker struct {
	Client *elastic.Client // elasticsearch client
	Index  string
}

func NewUnpacker(url, index string) unpacker {
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

	return unpacker{client, index}
}

func (u unpacker) PingElasticsearch(url string) {
	ctx := context.Background()
	info, code, err := u.Client.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	fmt.Printf("\nElasticsearch returned OK with code %d and version %s\n\n", code, info.Version.Number)
}

func (u unpacker) SetParsedStatus(status bool) {
	fmt.Println("All tags reset\n")

	termQuery := elastic.NewMatchAllQuery()
	ctx := context.Background()
	_, err := u.Client.UpdateByQuery(u.Index).
		Query(termQuery).
		Script(elastic.NewScript("ctx._source.parsed = params.status").Param("status", status)).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (u unpacker) SetFieldStringValue(fieldName string, s string) {

	termQuery := elastic.NewMatchAllQuery()
	ctx := context.Background()
	_, err := u.Client.UpdateByQuery(u.Index).
		Query(termQuery).
		Script(elastic.NewScript("ctx._source."+fieldName+" = params.data").Param("data", s)).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (u unpacker) SetFieldByteValue(fieldName string, b []byte) {

	termQuery := elastic.NewMatchAllQuery()
	ctx := context.Background()
	_, err := u.Client.UpdateByQuery(u.Index).
		Query(termQuery).
		Script(elastic.NewScript("ctx._source."+fieldName+" = params.data").Param("data", b)).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (u unpacker) ScritpByQuery(script string, params map[string]interface{}, query elastic.Query) {

	ctx := context.Background()
	_, err := u.Client.UpdateByQuery(u.Index).
		Query(query).
		Script(elastic.NewScript(script).Params(params)).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (u unpacker) ParseAndUpdate(fn func(map[string]interface{}) ([]byte, error)) {

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

	result, err := u.Client.Search().
		Index(u.Index).
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
		results := make(chan ParseResult, result.Hits.TotalHits)

		var wg sync.WaitGroup
		for _, hit := range result.Hits.Hits {
			jsonMap := make(map[string]interface{})

			// Deserialize hit.Source
			err := json.Unmarshal(*hit.Source, &jsonMap)
			if err != nil {
				// Deserialization failed
				fmt.Printf("Deserialization failed %v\n", err)
				continue
			}

			wg.Add(1)

			// Worker
			go func(id string) {
				defer wg.Done()

				json, err := fn(jsonMap)
				if err != nil {
					fmt.Println(err)
					return
				}
				results <- ParseResult{id, json}
			}(hit.Id)

		}

		// Closer
		go func() {
			wg.Wait()
			close(results)
		}()

		// Wait for parsing
		for res := range results {

			query := elastic.NewBoolQuery().Must(elastic.NewTermQuery("_id", res.documentId))
			script := "ctx._source.unpacked = params.data; ctx._source.parsed = true;"
			u.ScritpByQuery(script, map[string]interface{}{"data": res.data}, query)
			//u.SetFieldByQuery("parsed", true, query)
			//       fmt.Println("Parsepasrse\n", string(res.data))
		}

	} else {
		fmt.Print("\nFound no unparsed signals found\n")
	}
}

type ParseJob struct {
	documentId string
	json       map[string]interface{}
}

type ParseResult struct {
	documentId string
	data       []byte
}
