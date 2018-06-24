# ES-GUNPACK

Golang client that monitors one or multiple Elasticsearch indices and
1. Finds unparsed in ES-documents with binary data
2. Unpacks the data to JSON
3. Updates the documents with unpacked data

## Generate protobuf files
```
protoc --go_out=. pkg/example-protofiles/example.proto
```

## Generate go file for unpacking protofiles files
```
awk -f generate-struct-list.awk pkg/example-protofiles/*.pb.go
```


## Coming milestones
* 1 - The application that does a one time pass over the data in ES and parses value
* 2 - Protobuf signals can be added from user-defined directory
* 3 - The application continously monitors an index and parses as new documents comes in

## Todos 
* Parse protobuf signal based on signal name
* Add JSON to ElasticSearch
* Add JSON after parse
* Indicies as CLI arguments
* Password/user as optional input argument
* ...
* Decode ASN.1 messages (if needed)

## Done
* ~~Download and decode proto signal~
* ~~Compile protobuf signal~
* ~~Restructure project according to cli/pkg~~
* ~~Find a way of updating matched documents~~
* ~~Get data from Elasticsearch~~
* ~~Command line interface~~
* ~~Unmarshall Elasticsearch document to golang map~~
* ~~Update document with "parsed tag"~~


## Sprints
| Sprint # |                    Target                    |  Time (Actual)  |
|----------|----------------------------------------------|-----------------|
|          |                                              | 45 ()           |
| 8        | Learn about worker-pools                     | 45 (90)         |
| ~~7~~    | ~~Create goroutines that has the JSON data~~ | ~~45 ()~~       |
| 6        | Download and unpack protosignals             | 30 (42)         |
| 5        | Compile protobuf signal, upload data to ES   | 45 (57)         |
| 4        | Restructure project according to cli/pkg     | 45 (37)         |
| 3        | Update documents based on ID                 | 45 (71)         |
| 2        | Unmarshall json, ~~add parsed tag~~          | 45 (82)         |
| 1        | Use Elasticsearch client, retrieve data      | 45 (50)         |
|          |                                              | **Total**: 5h39 |

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
    - [Worker pools](https://gobyexample.com/worker-pools)
* Other
    - [Abstract Syntax Tree](https://zupzup.org/go-ast-traversal/) - parse Go files
    - AWK
        + [Online editor](https://www.tutorialspoint.com/execute_awk_online.php) 
        + [Good tutorial](https://www.ibm.com/developerworks/library/l-awk1/index.html)