package unpacker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

//
// Unpacker is a structure that holds data necessary to monitor and unpack data in Elasticsearch
//
type Unpacker struct {
	Client *elastic.Client // elasticsearch client
	Index  string
}

func (u Unpacker) SetParsedStatus(status bool) {
	fmt.Println("All tags reset\n")

	termQuery := elastic.NewMatchAllQuery()
	ctx := context.Background()
	_, err := u.Client.UpdateByQuery(u.Index).
		Query(termQuery).
		Script(elastic.NewScript("ctx._source.parsed = params.status").Param("status", status)).
		//    Script(elastic.NewScript("ctx._source.parsed = true").Params(map[string]interface{}{"tag": "blue"}).Lang("painless")).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (u Unpacker) ParseAndUpdate() {

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
			//    Script(elastic.NewScriptInline("ctx._source.parsed = false").Lang("painless")).
			//    Do(ctx)
		}
	} else {
		fmt.Print("\nFound no unparsed signals found\n")
	}
}
