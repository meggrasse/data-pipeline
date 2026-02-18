package main

import (
    // "fmt"
    "log"
    "grasse/pipeline"
	"os"
	"bufio"
	"strings"
)

type PressureTimeSeriesStream struct {}

// we should just share with the other, but whatever.
func (pressureTimeSeries *PressureTimeSeriesStream) Messages(stream chan pipeline.Message)  {
    file, err := os.Open("sample_data/source/pressure_timeseries_source.csv")
    if err != nil {
        // return nil, err
		// TODO: 
		log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	var count int

    for scanner.Scan() {
		row := scanner.Text()
        columns := strings.Split(row, ",")
        schemaType := columns[0]
        schemaVersion := columns[1]
        payload := strings.Join(columns[2:], ",")
		count++
        // TODO: use uuid
		stream <- pipeline.Message{Payload: payload, ID: count, SchemaType: schemaType, SchemaVersion: schemaVersion}
    }

    if err := scanner.Err(); err != nil {
		log.Fatal(err)
    }

	close(stream)

}