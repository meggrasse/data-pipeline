package main

import (
    "fmt"
    "log"
    // "grasse/pipeline"
	"os"
	"bufio"

)

type FlowTimeSeries struct {
	// I want to emulate streaming each from from sample_data/example_stream.csv.
	// Each column doesn't need to be typed; that can come at the data validation stage. It
	// should be raw.

	// stream chan string
}



// Reads line by line from the file specified by filename.
// Returns all lines as a slice of strings and an error if any.
// func (fts *FlowTimeSeries) ReadLines() ([]string, error) {
func (flowTimeSeries *FlowTimeSeries) Messages(stream chan string)  {
    // var lines []string
	// TODO:
    file, err := os.Open("sample_data/flow_timeseries.csv")
    if err != nil {
        // return nil, err
		// TODO: 
		log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		stream <- scanner.Text()
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

func main() {
	flowTimeSeries := FlowTimeSeries{}
	c := make(chan string)
	go flowTimeSeries.Messages(c)
	fmt.Println("start")
	for message := range c {
		fmt.Println(message)
	}
	fmt.Println("end")
}