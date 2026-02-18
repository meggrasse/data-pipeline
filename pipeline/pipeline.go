package pipeline

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MessageStream = chan Message

type Message struct {
	// TODO: bytes.
	Payload string
	ID uuid.UUID
	SchemaType string
	SchemaVersion string
}

type Source interface {
	// TODO: make sure it can only send.
	Messages(MessageStream)
}

type Processing interface {
	Process(MessageStream)
}

type Destination interface {
	// TODO: make sure it can only send.
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
	go merge(merged, channels)
	for _, processing := range p.Processings {
		go processing.Process(merged)
	}
	p.Destination.Messages(merged)
}

// merge fans in from all source channels; validates schema/version and forwards only valid messages.
// Invalid messages are skipped (clients handle their own errors in each stage).
func merge(out MessageStream, inputs []MessageStream) {
	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(c MessageStream) {
			for message := range c {
				versions, ok := Schemas[message.SchemaType]
				if !ok {
					fmt.Printf("Skipping message %s: schema not supported: %s. Allowed: %v\n", message.ID, message.SchemaType, Schemas)
					continue
				}
				valid := false
				for _, version := range versions {
					if version == message.SchemaVersion {
						valid = true
						break
					}
				}
				if !valid {
					fmt.Printf("Skipping message %s: schema version not supported: %s. Allowed: %v\n", message.ID, message.SchemaVersion, versions)
					continue
				}
				out <- message
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

var Schemas = map[string][]string{
	"timeseries.flow": {"1.0.0"},
	"timeseries.pressure": {"1.0.0"},
}

// todo: send and receive channels only in source/destination
// log
// zero or more processing stages that
// transform or validate data,

// some processing stages transform (modify, reduce filter) messages based on a single schema.
// others may combine multiple messages (i.e., with the same timestamp and site internal id) into one representation. 

// channels from all sources need to be available to the processesing stage.
	// Or actually from previous processing stages as well.
	// Messages themselves likely need an indentifer, and then sources need to
	// describe themselves and their schema.
	// New schemas may be defined by processing stages as they provide transformations to later stages.