# Useful

## fmt.Printf("%+v\n", result)
%+v prints struct keys and values:
```
&{TookInMillis:3 ScrollId: Hits:0xc4201e6240 Suggest:map[] Aggregations:map[] TimedOut:false Error:<nil> Profile:<nil> Shards:0xc4201e8060}
```


# Postman 
### Postman community collections
There are several useful collections that can be loaded into Postman. I found this one for Elasticsearch:
https://github.com/sittinash/elasticsearch-postman/blob/master/Elasticsearch%20APIs.postman_collection.json

**The examples overrides my own defined index name**. Can be changed under Pre-request Script
```
postman.setEnvironmentVariable("index", "logstash-2017.10.24");
```

### Postman variables
Variables for insertion into requests can be defined in Postman, for example:
```
url:127.0.0.1
port:9200
index:logstash-2018.06.15
```

url:127.0.0.1
port:9200
index:logstash-2018.06.15

### Postman Authorization
You can insert various types of authentication into your requests. Available 
I added the username/password needed for authentication to Elasticsearch under Get mapping -> Authorization -> Type (Basic Auth) -> enter username/password.

### Short ones
* [Kibana dev console](http://localhost:5601/app/kibana#/dev_tools/console?_g=()) is useful to debug queries and filters
* [Boolean filters](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-bool-query.html) are used for boolean logic in Elasticsearch queries
```
POST _search
{
  "query": {
    "bool" : {
      "must" : {
        "term" : { "user" : "kimchy" }
      },
      "filter": {
        "term" : { "tag" : "tech" }
      },
      "must_not" : {
        "range" : {
          "age" : { "gte" : 10, "lte" : 20 }
        }
      },
      "boost" : 1.0
    }
  }
}
```
A general guide is ([SO: elasticsearch bool query combine must with OR](https://stackoverflow.com/questions/28538760/elasticsearch-bool-query-combine-must-with-or))
```
OR is spelled "should"
AND is spelled "must"
NOR is spelled "should_not"
```

# Bugs/weird things

# Elastic-lib
Really a feature of Elasticsearch, but when using Painless scripting language, the variable has to be called as **params.status**. E.g.:
**Correct:**
```
Script(elastic.NewScript("ctx._source.parsed = params.status").Param("status", status)).
```
**Incorrect:**
```
Script(elastic.NewScript("ctx._source.parsed = status").Param("status", status)).
```
```
elasticsearch_1  |  at java.lang.Thread.run(Thread.java:748) [?:1.8.0_141]
elasticsearch_1  | Caused by: java.lang.IllegalArgumentException: Variable [status] is not defined.
elasticsearch_1  |  at org.elasticsearch.painless.PainlessScript$Script.compile(ctx._source.parsed = status:22) ~[?:?]
```


## Can't update multiple values in single document
Using the Update() function, only the last "Script()" is actually triggered. Below, only ***_source.something_else*** will be updated, and _source.parsed left as it is.
```
update, _ := client.Update().Index(updateIndex).Type("doc").Id("AWQCkaaoxTNcL8eIfO44").
    Script(elastic.NewScriptInline("ctx._source.parsed = false").Lang("painless")).
    Script(elastic.NewScriptInline("ctx._source.something_else = false").Lang("painless")).
    Do(ctx)

```

## elastic.SetSniff(false)
Initially, I was unable to connect to my local Elasticsearch instance with the client library. The connection failed with the message:
```
2018/06/16 20:32:36 no active connection found: no Elasticsearch node available
```

In the end I found the answer in the Elktail application: https://github.com/knes1/elktail/blob/master/elktail.go . It seems that setting the elastic.SetSniff(false), option when creating the client solves the issue. I don't know why.

    defaultOptions := []elastic.ClientOptionFunc{
        elastic.SetURL("http://127.0.0.1:9200", "http://localhost:9200"),
        elastic.SetSniff(false),
        elastic.SetBasicAuth("elastic", "changeme"),
        elastic.SetHealthcheckTimeoutStartup(10 * time.Second),
        elastic.SetHealthcheckTimeout(2 * time.Second),
    }

## Print query
Useful for debugging
```
src, err := boolTermQuery.Source()
if err != nil {
    panic(err)
}
data, err := json.MarshalIndent(src, "", "  ")
if err != nil {
    panic(err)
}
fmt.Println(string(data))
```
Result:
```
{
  "bool": {
    "must": {
      "term": {
        "parsed": true
      }
    }
  }
}
```

# Install protoc-gen-go

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

# Sublime

| Command | Description |
|---------|-------------|
| ctrl+g  | gotoline    |