package main

import (
    "grasse/pipeline"
)

type EmptySource struct {}

var EmptySchemaType = "empty"
var EmptySchemaVersion = "1.0.0"

func (emptySource *EmptySource) Messages(stream chan pipeline.Message)  {
	stream <- pipeline.Message{Payload: "empty", ID: 0, SchemaType: EmptySchemaType, SchemaVersion: EmptySchemaVersion}
	close(stream)
}