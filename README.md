This was a take-home project for an interview process. Buyer beware!

# Execution

## Running the example

To run the provided `example` module that exercises the `pipeline` package, navigate to the `example` directory and run `$ go run .`. The output from the example will be populated in `example/sample_data/destination/sensor_timeseries_destination.csv`. The example includes input files, source stages, processing stages, and a destination stage that writes the stream to the output file.

## Tests

Tests are included in the `pipeline` package. Those can be exercised from the `pipeline` directory using `$ go test`.

# Design

## Design Goals

- Stream messages from multiple sources simultaneously.
- Validate messages against schemas known by the pipeline (this is intentionally naive given the fidelity of this system; see [Future Improvements](#future-improvements))
- Allow for transformation via 0 or more serial processing stages.
- Program should exit once all source channels are closed.
- Aggregation (i.e. collapsing multiple messages into one) in any stage should be supported.

## Overall architecture

The framework provides a `Pipeline` struct, and `Source`, `Processing` and `Destination` interfaces to provide clients some scaffolding in building their data pipelines. 

`Pipeline` supports the following functionality:

- `Message`s can be streamed from one or more sources, and will be fanned in by the pipeline into a single `MessageStream` in the order they were received.
- Whilst fanning in, the pipeline ensures the `Message`'s schema/version pair is valid (i.e., known by the pipeline).
- Processing stages are described as a collection and must be serial. The expected practice from the client is to check the schema of the incoming message, and if it doesn't process messages of this schema, pass it through to the `out` stream.
- Source and processing stages are responsible for closing `out` once they've drained their input. This matches the Go idiom where producers are responsible for closing channels (the stage is considered the producer here as they are the only writer to the channel, even though they don't technically create it). 
- A single destination receives the `MessageStream` (the `out` stream in the final processing stage) to route to a destination (i.e., write to a file). It's important that destination implementations don't close the incoming channel after processing, as producers close.

## Concurrency Model

Each source and processing stage are executed via goroutines to allow messages to be streamed through the system (if all stages were executed on the main goroutine, all data and `close` would need to be sent over the channel before the next stage could begin). The destination stream runs on the main goroutine to ensure the program doesn't terminate until that message stream finishes executing.

Source messages are additionally streamed via distinct goroutines to allow for concurrent streaming.

To ensure graceful termination (exiting the program once all source channels are closed), we need to ensure the fan-in message stream is closed once all the source channels are closed. This is achieved using the Go idiom of a channel as a semaphore: the channel is sent a value once each source stream is closed, waits on receiving a value for each channel, and then closes `out`.

## Testing/Validation

Units tests cover boundary conditions (e.g., no source stages, no processing stages), the schema validation contract, and a simple happy-path pipeline.

The example adoption of course also helps validate the pipeline.

## Limitations

- None of the channels are buffered and all risk running out of system resources (causing program to exit). If this system was to be deployed, a reasonable buffer length should be determined via memory availability on target hardware and empirical use.
- Lots of client responsibility: Processing stages check each incoming message's schema and pass message through if not applicable, source and processing stage close in channels, and there's no error handling infrastructure in the pipeline -- stages are responsible for locally handling their own errors. 
- Message data are strings, so data representing a series of bytes must be encoded as base64 (or similar).
- Keys in `Pipeline.Schemas` should really be renamed to express they are based on CSV.

## Future Improvements

- Processing interface includes schema types/versions supported by that stage. System is responsible for validating messages into each stage, and passing them directly to the out stream so clients don't have to own this logic, and also allows for stronger validation between processing stages. This also means supported schema types/versions are dynamically set by clients rather than hardcoded in `Pipeline`.
