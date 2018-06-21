# ES-GUNPACK

Golang client that monitors one or multiple Elasticsearch indices and
1. Finds unparsed in ES-documents with binary data
2. Unpacks the data to JSON
3. Updates the documents with unpacked data

## Generate protobuf files
```
protoc --go_out=. pkg/example-protofiles/example.proto
```

## Coming milestones
* 1 - The application that does a one time pass over the data in ES and parses value
* 2 - Protobuf signals can be added from user-defined directory
* 3 - The application continously monitors an index and parses as new documents comes in

## Todos 
* Parse protobuf signal based on signal name
* Download and decode proto signal
* Add JSON to ElasticSearch
* Add JSON after parse
* Indicies as CLI arguments
* Password/user as optional input argument
* ...
* Decode ASN.1 messages (if needed)

## Done
* ~~Compile protobuf signal~
* ~~Restructure project according to cli/pkg~~
* ~~Find a way of updating matched documents~~
* ~~Get data from Elasticsearch~~
* ~~Command line interface~~
* ~~Unmarshall Elasticsearch document to golang map~~
* ~~Update document with "parsed tag"~~


## Sprints
| Sprint # |                   Target                   |  Time (Actual)  |
|----------|--------------------------------------------|-----------------|
|          |                                            | 45 ()           |
|        6 | Download and unpack protosignals           | 30 (42)           |
|        5 | Compile protobuf signal, upload data to ES | 45 (57)         |
|        4 | Restructure project according to cli/pkg   | 45 (37)         |
|        3 | Update documents based on ID               | 45 (71)         |
|        2 | Unmarshall json, ~~add parsed tag~~        | 45 (82)         |
|        1 | Use Elasticsearch client, retrieve data    | 45 (50)         |
|          |                                            | **Total**: 5h29 |

##
* Elasticsearch in general
    - [Query vs Filter](https://www.elastic.co/blog/found-optimizing-elasticsearch-searches) - 
    - [BoolQuery, Must, Must_not](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-bool-query.html)
    - [Painless scripting](https://www.elastic.co/guide/en/elasticsearch/painless/current/painless-examples.html)
* elasticsearch client
    - https://godoc.org/github.com/olivere/elastic
    - https://github.com/olivere/elastic/wiki
* cli refererence - https://github.com/urfave/cli#generated-help-text
* Go best practices
    - [Peter Bourgon - Best practices 2016](https://peter.bourgon.org/go-best-practices-2016/#dependency-management)
    - [Template for repository structure](https://github.com/thockin/go-build-template)