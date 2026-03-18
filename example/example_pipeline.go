package main

import (
	"grasse/pipeline"
)

func main() {
	flowSource := SensorTimeSeriesStream{filename: "flow_timeseries_source.csv"}
	pressureSource := SensorTimeSeriesStream{filename: "pressure_timeseries_source.csv"}
	emptySource := EmptySource{}
	destination := SensorTimeSeriesDestination{filename: "sensor_timeseries_destination.csv"}
	transformation := ExampleTransformation{}
	simpeTransformation := SimpleTransformation{}

	pipeline := pipeline.Pipeline{
		Sources: []pipeline.Source{&flowSource, &pressureSource, &emptySource},
		Processings: []pipeline.Processing{&transformation, &simpeTransformation},
		Destination: &destination,
	}
	pipeline.Run()
}