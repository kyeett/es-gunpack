# ES-GUNPACK

Golang client that monitors one or multiple Elasticsearch indices and
1. Finds unparsed in ES-documents with binary data
2. Unpacks the data to JSON
3. Updates the documents with unpacked data



## Todos
* Unmarshall Elasticsearch document to golang map
* Update document with "parsed tag"

## Done
* ~~Get data from Elasticsearch~~
* ~~Command line interface~~

## Sprints
| Sprint # |                  Target                 |  Time (Actual)  |
|----------|-----------------------------------------|-----------------|
|        2 | Unmarshall json, ~~add parsed tag~~     | 45 (82)         |
|        1 | Use Elasticsearch client, retrieve data | 45 (50)         |
|          |                                         | **Total**: 0h45 |

##
* elasticsearch client - https://godoc.org/github.com/olivere/elastic
* cli refererence - https://github.com/urfave/cli#generated-help-text