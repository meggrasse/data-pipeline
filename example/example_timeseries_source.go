package main

import (
  "grasse/pipeline"
	"os"
	"bufio"
	"strings"
  "github.com/google/uuid"
)

type SensorTimeSeriesStream struct {
  filename string
}

func (ts *SensorTimeSeriesStream) Messages(c chan pipeline.Message) {
    file, open_err := os.Open("sample_data/source/" + ts.filename)
    if open_err != nil {
      return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	  var count int

    for scanner.Scan() {
      row := scanner.Text()
      fields := strings.Split(row, ",")

      // Extra field validation should go here.
      // For now we'll assume we know the shape of the data since we own the source
      // and first two fields define the schema.
      schemaType := fields[0]
      schemaVersion := fields[1]
      payload := strings.Join(fields[2:], ",")
      count++
      c <- pipeline.Message{Payload: payload, ID: uuid.New(), SchemaType: schemaType, SchemaVersion: schemaVersion}
    }

    if err := scanner.Err(); err != nil {
    }

    close(c)
}