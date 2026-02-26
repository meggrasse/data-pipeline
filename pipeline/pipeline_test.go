package pipeline

import (
	"testing"
	"github.com/google/uuid"
)

func TestSupportedSchema(t *testing.T) {
	var validSchemaType string
	var validSchemaVersion string
	// TODO: make sure schemas isn't empty
	for schemaType, versions := range Schemas {
		validSchemaType = schemaType
		validSchemaVersion = versions[0]
		break
	}
	message := Message {
		Payload: "test",
		ID: uuid.New(),
		SchemaType: validSchemaType,
		SchemaVersion: validSchemaVersion,
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go merge(out, []MessageStream{in})
	in <- message
	close(in)
	count := 0
	for _ = range out {
		count++
	}
	if count != 1 {
		t.Errorf("Message should be valid: %d", count)
	}
}

func TestUnsupportedSchemaType(t *testing.T) {
	message := Message{
		Payload: "test",
		ID: uuid.New(),
		SchemaType: "not-supported",
		SchemaVersion: "1.0.0",
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go merge(out, []MessageStream{in})
	in <- message
	close(in)
	for _ = range out {
		t.Errorf("Message should not be valid")
	}
}

func TestUnsupportedSchemaVersion(t *testing.T) {
	var validSchemaType string
	// TODO: make sure schemas isn't empty
	for schemaType, _ := range Schemas {
		validSchemaType = schemaType
		break
	}
	message := Message {
		Payload: "test",
		ID: uuid.New(),
		SchemaType: validSchemaType,
		SchemaVersion: "not-supported",
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go merge(out, []MessageStream{in})
	in <- message
	close(in)
	for _ = range out {
		t.Errorf("Message should not be valid")
	}
}

func testMessage() Message {
	var validSchemaType string
	var validSchemaVersion string
	// TODO: make sure schemas isn't empty
	for schemaType, versions := range Schemas {
		validSchemaType = schemaType
		validSchemaVersion = versions[0]
		break
	}
	return Message{
		Payload: "test",
		ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		SchemaType: validSchemaType,
		SchemaVersion: validSchemaVersion,
	}
}

type TestSource struct {
	t *testing.T
}

func (testSource *TestSource) Messages(stream chan Message) {
	stream <- testMessage()
	close(stream)
}

type TestDestination struct {
	t *testing.T
}

func (testDestination *TestDestination) Messages(stream chan Message) {
	count := 0
	for message := range stream {
		count++
		if message != testMessage() {
			testDestination.t.Errorf("Message should be unaltered: %v", message)
		}
	}
	if count != 1 {
		testDestination.t.Errorf("Expected single message.")
	}
}

// TODO: test no source?
func TestNoProcessingPipeline(t *testing.T) {

	pipeline := Pipeline {
		Sources: []Source{&TestSource{t: t}},
		Processings: []Processing{},
		Destination: &TestDestination{t: t},
	}
	pipeline.Run()
}