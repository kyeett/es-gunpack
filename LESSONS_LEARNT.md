# Useful

## fmt.Printf("%+v\n", result)
%+v prints struct keys and values:
```&{TookInMillis:3 ScrollId: Hits:0xc4201e6240 Suggest:map[] Aggregations:map[] TimedOut:false Error:<nil> Profile:<nil> Shards:0xc4201e8060}
```

# Bugs/weird things

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


