package main

import (
	"grasse/pipeline"
)

func main() {
	flowSource := SensorTimeSeriesStream{filename: "flow_timeseries_source.csv"}
	pressureSource := SensorTimeSeriesStream{filename: "pressure_timeseries_source.csv"}
	emptySource := EmptySource{}
	destination := SensorTimeSeriesDestination{filename: "sensor_timeseries_destination.csv"}
	aggregation := ExampleAggregation{}
	unitConversion := UnitConversion{}

	pipeline := pipeline.Pipeline{
		Sources: []pipeline.Source{&flowSource, &pressureSource, &emptySource},
		Processings: []pipeline.Processing{&aggregation, &unitConversion},
		Destination: &destination,
	}
	pipeline.Run()
}