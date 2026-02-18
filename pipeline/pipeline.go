package pipeline

import (
	"fmt"
	"sync"
)

type Pipeline struct {
	Sources []Source
	Processings []Processing
	// Single destination only
	Destination Destination
}

func (p *Pipeline) Run() {
	var channels []chan Message
	for _, source := range p.Sources {
		ch := make(chan Message)
		channels = append(channels, ch)
		go source.Messages(ch)
	}
	merged := make(chan Message)
	go merge(merged, channels)
	for _, processing := range p.Processings {
		go processing.Process(merged)
	}
	p.Destination.Messages(merged)
}

func merge(out chan Message, inputs []chan Message) {
	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(c <-chan Message) {
			for message := range c {
				versions, ok := Schemas[message.SchemaType]
				if !ok {
					fmt.Println("Skipping message: %s schema not supported: %s. Allowed schemas: %v", message.ID, message.SchemaType, Schemas)
					continue
				}
				valid := false
				for _, v := range versions {
					if v == message.SchemaVersion {
						valid = true
						break
					}
				}
				if !valid {
					fmt.Println("Skipping message: %s schema version not supported: %s. Allowed versions: %v", message.ID, message.SchemaVersion, versions)
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
// result (i.e. errors)
// log
// zero or more processing stages that
// transform or validate data,
