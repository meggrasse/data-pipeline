package main

import (
	"fmt"
	"grasse/pipeline"
	"strings"
	"github.com/google/uuid"
	"strconv"
)

type ExampleAggregation struct {}

func (agg *ExampleAggregation) Process(in pipeline.MessageStream, out pipeline.MessageStream) {
	defer close(out)
	var currentFlowTimestamp string
	var currentPressureTimestamp string
	var flowData float64
	var pressureData float64
	var currentData float64
	for message := range in {
		if message.SchemaType != "timeseries.flow" && message.SchemaType != "timeseries.pressure" || message.SchemaVersion != "1.0.0" {
			out <- message
			continue
		}
		fields := strings.Split(message.Payload, ",")
		// Assuming shape per schema.
		timestamp := fields[0]
		currentData, _ = strconv.ParseFloat(fields[len(fields)-1], 64)
		if message.SchemaType == "timeseries.flow" {
			if currentFlowTimestamp == timestamp {
				flowData += currentData
			} else {
				out <- pipeline.Message{
					Payload: fmt.Sprintf("flow,%s,%f", currentFlowTimestamp, flowData),
					ID: uuid.New(),
					SchemaType: "timeseries.batched",
					SchemaVersion: "1.0.0",
				}
				flowData = currentData
				currentFlowTimestamp = timestamp
			}
		} else if message.SchemaType == "timeseries.pressure" {
			if currentPressureTimestamp == timestamp {
				pressureData += currentData
			} else {
				out <- pipeline.Message{
					Payload: fmt.Sprintf("pressure,%s,%f", currentPressureTimestamp, pressureData),
					ID: uuid.New(),
					SchemaType: "timeseries.batched",
					SchemaVersion: "1.0.0",
				}
				pressureData = currentData
				currentPressureTimestamp = timestamp
			}
		}
	}
}

type UnitConversion struct {}

func (conv *UnitConversion) Process(in pipeline.MessageStream, out pipeline.MessageStream) {
	defer close(out)
	for message := range in {
		if message.SchemaType != "timeseries.batched" || message.SchemaVersion != "1.0.0" {
			out <- message
			continue
		}
		fields := strings.Split(message.Payload, ",")
		var unit string
		// Assuming shape per schema.
		kind := fields[0]
		timestamp := fields[1]
		// Since we're not modifying it, just keep `currentData` as a string.
		currentData := fields[2]
		switch kind {
			case "flow":
				unit = "m3/s"
			case "pressure":
				unit = "bar"
			default:
				out <- message
				continue
			}
		payload := fmt.Sprintf("%s,%s,%s,%s", kind, unit, timestamp, currentData)  	
		out <- pipeline.Message{
			Payload: payload,
			ID: uuid.New(),
			SchemaType: "timeseries.batched",
			SchemaVersion: "1.1.0",
		}
	}
}
