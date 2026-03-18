package pipeline

import (
	"log"
	"maps"
	"slices"

	"github.com/google/uuid"
)

type MessageStream = chan Message

type Message struct {
	Payload string
	ID uuid.UUID
	SchemaType string
	SchemaVersion string
}

type Source interface {
	Messages(MessageStream)
}

type Processing interface {
	Process(in MessageStream, out MessageStream)
}

type Destination interface {
	Messages(MessageStream)
}

type Pipeline struct {
	Sources []Source
	Processings []Processing
	// Single destination only
	Destination Destination
}

func (p *Pipeline) Run() {
	var channels []MessageStream
	for _, source := range p.Sources {
		ch := make(MessageStream)
		channels = append(channels, ch)
		go source.Messages(ch)
	}
	merged := make(MessageStream)
	go fanIn(merged, channels)
	ch := merged
	for _, proc := range p.Processings {
		next := make(MessageStream)
		go proc.Process(ch, next)
		ch = next
	}
	p.Destination.Messages(ch)
}

// Blocks until all sources are closed; should be called in a goroutine.
func fanIn(out MessageStream, ins []MessageStream) {
	// Following Go's idiom of using a channel as a semaphore.
	closed := make(chan int, len(ins))
	for _, ch := range ins {
		// Unique loop var on each iteration as of 1.22.
		// c.f. https://go.dev/wiki/LoopvarExperiment 
		go func() {
			defer func() { closed <- 1 }()
			for msg := range ch {
				if validate(msg) {
					out <- msg
				}
			}
		}()
	}
	for range len(ins) {
		<- closed
	}
	close(out)
}

func validate(msg Message) bool {
	versions, ok := Schemas[msg.SchemaType]
	if !ok {
		log.Printf("Message not valid: schema type not supported [Message ID: %s; Schema type: %s] Supported Schema types: %v \n", msg.ID, msg.SchemaType, slices.Collect(maps.Keys(Schemas)))
		return false
	}
	for _, version := range versions {
		if version == msg.SchemaVersion {
			return true
		}
	}
	log.Printf("Message not valid: schema version not supported [Message ID: %s; Schema version: %s] Supported Schema versions: %v \n", msg.ID, msg.SchemaVersion, versions)
	return false
}

var Schemas = map[string][]string{
	"timeseries.flow": {"1.0.0"},
	"timeseries.pressure": {"1.0.0"},
	"timeseries.batched": {"1.0.0", "1.1.0"},
}
