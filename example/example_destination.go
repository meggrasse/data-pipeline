package main

import (
	"os"
	"log"
	"grasse/pipeline"
)

type FlowTimeSeriesDestination struct {
	// I want to emulate writing to a file specified by filename.
	// Each line should be a string.

	// stream chan string
}

func (flowTimeSeries *FlowTimeSeriesDestination) Messages(stream chan pipeline.Message) {
	// TODO:
	file, err := os.Create("sample_data/flow_timeseries_destination.csv")
	if err != nil {
		// TODO:
		log.Fatal(err)
	}
	defer file.Close()

	for message := range stream {
		file.WriteString("[" + message.SchemaType + "," + message.SchemaVersion + "] - (" + message.Payload + ") \n")
	}
}