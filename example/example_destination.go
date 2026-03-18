package main

import (
	"os"
	"log"
	"grasse/pipeline"
)

type SensorTimeSeriesDestination struct {
  filename string
}

func (ts *SensorTimeSeriesDestination) Messages(c chan pipeline.Message) {
	file, err := os.Create("sample_data/destination/" + ts.filename)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	for message := range c {
		file.WriteString(message.ID.String() + ": [" + message.SchemaType + "," + message.SchemaVersion + "] - (" + message.Payload + ") \n")
	}
}