package pipeline

import (
	"fmt"
	"sync"

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
	go merge(merged, channels)
	ch := merged
	for _, proc := range p.Processings {
		next := make(MessageStream)
		go proc.Process(ch, next)
		ch = next
	}
	p.Destination.Messages(ch)
}

// TODO use wg.GO
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
	"timeseries.batched": {"1.0.0", "1.1.0"},
}
