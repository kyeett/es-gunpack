package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	example "github.com/kyeett/es-gunpack/pkg/example-protofiles"
	"github.com/kyeett/es-gunpack/pkg/unpacker"
)

// Simple cli that adds a protobuf signal to the 'data' field of all entries in logstash-all
func main() {
	test := &example.Test{
		Label: proto.String("hello proto!"),
		Type:  proto.Int32(17),
	}
	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	dataTextString := proto.MarshalTextString(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	newTest := &example.Test{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if test.GetLabel() != newTest.GetLabel() {
		log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
	}
	// etc.
	fmt.Printf("%+v, %T\n", data, data)
	fmt.Printf("%+v, %T\n", dataTextString, dataTextString)

	// Create client interfacing elasticsearch
	url := "http://localhost:9200"
	unpackerClient := unpacker.NewUnpacker(url, "logstash-2018.06.15")

	// Update payload
	unpackerClient.SetFieldByteValue("data", data)

}
