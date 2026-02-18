package main

import (
	"grasse/pipeline"
)

func main() {
	flowSource := FlowTimeSeriesStream{}
	pressureSource := PressureTimeSeriesStream{}
	emptySource := EmptySource{}
	destination := FlowTimeSeriesDestination{}
	transformation := ExampleTransformation{}
	pipeline := pipeline.Pipeline{Sources: []pipeline.Source{&flowSource, &pressureSource, &emptySource}, Destination: &destination, Processings: []pipeline.Processing{&transformation}}
	pipeline.Run()
}