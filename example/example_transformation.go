package main

import "grasse/pipeline"

type ExampleTransformation struct {}

func (exampleTransformation *ExampleTransformation) Process(channel chan pipeline.Message) {
	for message := range channel {
		if message.SchemaType == "timeseries.flow" && message.SchemaVersion == "1.0.0" {
			channel <- pipeline.Message{Payload: message.Payload + "flow_timeseries transformed", ID: message.ID, SchemaType: message.SchemaType, SchemaVersion: message.SchemaVersion}
		} else if message.SchemaType == "timeseries.pressure" && message.SchemaVersion == "1.0.0" {
			channel <- pipeline.Message{Payload: message.Payload + "pressure_timeseries transformed", ID: message.ID, SchemaType: message.SchemaType, SchemaVersion: message.SchemaVersion}
		}
	}
}