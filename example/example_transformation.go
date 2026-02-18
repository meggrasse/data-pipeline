package main

import "grasse/pipeline"

type ExampleTransformation struct {
}

// i would make a map of channels using schema type and version as the key but there might be duplicates.
// it's actually probably reasonable to assume a transformation operates on a single schema
func (exampleTransformation *ExampleTransformation) Process(channel chan pipeline.Message) {
	// we assume this transformation is only for a specific schema.
	for message := range channel {
		if message.SchemaType == "timeseries.flow" && message.SchemaVersion == "1.0.0" {
			channel <- pipeline.Message{Payload: message.Payload + "flow_timeseries transformed", ID: message.ID, SchemaType: message.SchemaType, SchemaVersion: message.SchemaVersion}
		} else if message.SchemaType == "timeseries.pressure" && message.SchemaVersion == "1.0.0" {
			channel <- pipeline.Message{Payload: message.Payload + "pressure_timeseries transformed", ID: message.ID, SchemaType: message.SchemaType, SchemaVersion: message.SchemaVersion}
		}
	}
}