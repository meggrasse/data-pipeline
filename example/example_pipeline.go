package main

import (
	"grasse/pipeline"
)

func main() {
	source := FlowTimeSeriesStream{}
	destination := FlowTimeSeriesDestination{}
	pipeline := pipeline.Pipeline{Sources: []pipeline.Source{&source}, Destination: &destination}
	pipeline.Run()
}