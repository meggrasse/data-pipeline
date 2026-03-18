package pipeline

import (
	"testing"
	"github.com/google/uuid"
	"sort"
)

type TestSource struct {}

type TestProcessing struct {}

type TestDestination struct {
	t *testing.T
}

type TestNoSourceDestination struct {
	t *testing.T
}

func (dest *TestNoSourceDestination) Messages(stream MessageStream) {
	_, open := <- stream
	if open {
		dest.t.Errorf("Expected stream to be closed.")
	}
}

func validSchema() (string, string) {
	var keys []string
	for k := range Schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys[0], Schemas[keys[0]][0]
}

func testMessage() Message {
	validSchemaType, validSchemaVersion := validSchema()
	return Message{
		Payload: "test",
		ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		SchemaType: validSchemaType,
		SchemaVersion: validSchemaVersion,
	}
}

func (testSource *TestSource) Messages(stream MessageStream) {
	stream <- testMessage()
	close(stream)
}

func (tp *TestProcessing) Process(in MessageStream, out MessageStream) {
	for message := range in {
		out <- message
	}
	close(out)
}

func (testDestination *TestDestination) Messages(stream MessageStream) {
	count := 0
	for message := range stream {
		count++
		msg := testMessage()
		if message != msg {
			testDestination.t.Errorf("Message should be unaltered: %v != %v", message, msg)
		}
	}
	if count != 1 {
		testDestination.t.Errorf("Expected single message.")
	}
}

func TestSupportedSchema(t *testing.T) {
	validSchemaType, validSchemaVersion := validSchema()
	message := Message {
		Payload: "test",
		ID: uuid.New(),
		SchemaType: validSchemaType,
		SchemaVersion: validSchemaVersion,
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go fanIn(out, []MessageStream{in})
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
		SchemaType: "not-supported-type",
		SchemaVersion: "1.0.0",
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go fanIn(out, []MessageStream{in})
	in <- message
	close(in)
	for _ = range out {
		t.Errorf("Message should not be valid")
	}
}

func TestUnsupportedSchemaVersion(t *testing.T) {
	validSchemaType, _ := validSchema()
	message := Message {
		Payload: "test",
		ID: uuid.New(),
		SchemaType: validSchemaType,
		SchemaVersion: "not-supported-vers",
	}
	in := make(MessageStream)
	out := make(MessageStream)
	go fanIn(out, []MessageStream{in})
	in <- message
	close(in)
	for _ = range out {
		t.Errorf("Message should not be valid")
	}
}

func TestPipeline(t *testing.T) {
	pipeline := Pipeline {
		Sources: []Source{&TestSource{}},
		Processings: []Processing{&TestProcessing{}},
		Destination: &TestDestination{t: t},
	}
	pipeline.Run()
}

func TestNoSourcePipeline(t *testing.T) {
	pipeline := Pipeline {
		Sources: []Source{},
		Processings: []Processing{&TestProcessing{}},
		Destination: &TestNoSourceDestination{t: t},
	}
	pipeline.Run()
}

func TestNoProcessingPipeline(t *testing.T) {
	pipeline := Pipeline {
		Sources: []Source{&TestSource{}},
		Processings: []Processing{},
		Destination: &TestDestination{t: t},
	}
	pipeline.Run()
}

func TestFanIn(t *testing.T) {
	validSchemaType, validSchemaVersion := validSchema()
	m := Message{
		Payload: "test",
		ID: uuid.New(),
		SchemaType: validSchemaType,
		SchemaVersion: validSchemaVersion,
	}
	in1 := make(MessageStream)
	in2 := make(MessageStream)
	out := make(MessageStream)
	go fanIn(out, []MessageStream{in1, in2})
	go func() {
		in1 <- m
		in2 <- m
	}()
	// drain out values whilst both channels are still open
	<- out
	<- out

	// ensure `close` is sent on `out` if input channels are closed.
	close(in1)
	close(in2)
	_, open := <- out
	if open {
		t.Errorf("Expected out channel to be closed")
	}
}