package main

import (
    "grasse/pipeline"
	"github.com/google/uuid"
)

// This exists to exemplify a source with an unsupported schema.

type EmptySource struct {}

func (emptySource *EmptySource) Messages(stream chan pipeline.Message)  {
	const schemaType = "empty"
	const schemaVersion = "1.0.0"
	stream <- pipeline.Message{Payload: "empty", ID: uuid.New(), SchemaType: schemaType, SchemaVersion: schemaVersion}
	close(stream)
}