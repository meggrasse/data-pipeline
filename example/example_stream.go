package main

import (
    // "fmt"
    "log"
    "grasse/pipeline"
	"os"
	"bufio"

)

type FlowTimeSeriesStream struct {
	// I want to emulate streaming each from from sample_data/example_stream.csv.
	// Each column doesn't need to be typed; that can come at the data validation stage. It
	// should be raw.

	// stream chan string
}

// TODO: extract from data
var SchemaType = "flow_timeseries"
var SchemaVersion = "1.0.0"

// Reads line by line from the file specified by filename.
// Returns all lines as a slice of strings and an error if any.
// func (fts *FlowTimeSeries) ReadLines() ([]string, error) {
func (flowTimeSeries *FlowTimeSeriesStream) Messages(stream chan pipeline.Message)  {
    // var lines []string
	// TODO:
    file, err := os.Open("sample_data/flow_timeseries_source.csv")
    if err != nil {
        // return nil, err
		// TODO: 
		log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	var count int

    for scanner.Scan() {
		count++
		stream <- pipeline.Message{Payload: scanner.Text(), ID: count, SchemaType: SchemaType, SchemaVersion: SchemaVersion}
        // lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        // return nil, err
		// TODO:
		log.Fatal(err)
    }

	close(stream)
    // return lines, nil
}

// func main() {
// 	flowTimeSeries := FlowTimeSeriesStream{}
// 	c := make(chan string)
// 	go flowTimeSeries.Messages(c)
// 	fmt.Println("start")
// 	for message := range c {
// 		fmt.Println(message)
// 	}
// 	fmt.Println("end")
// }