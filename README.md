This was a take-home project for an interview process. Buyer beware!

## Design Goals

- Support multiple, simultanous sources.
- Allow for window processing.

## Overall architecture

The framework provides a `Pipeline` struct, and `Source`, `Processing` and `Destination` interfaces to provide clients some scaffolding in building their data pipelines. 

`Pipeline` supports the following functionality:

    - `Message`s can be streamed one of more sources, and will be merged together by the pipeline into a single `MessageStream` in the order they were receieved. Merging is acheived via Go's `WaitGroup` (i.e., a semaphore).
    - Whilst merging, the pipeline ensures the the `Message`'s schema/version pair is valid (i.e., known by the pipeline). TODO: doesn't the client do the proessing?
    - Processing stages are described as a collection and must be linear. The expected practice from the client is to check the schema of the incoming message, and if it doesn't process messages of this schema, pass it through to the `out` stream. TODO: this suggests that processing stages could actually be defined alongside known schemas by clients. And then better validation checking could be done by the pipeline.
    - A single destination recieves the `MessageStream` (the `out` stream in the final processing stage) to route to a destination (i.e., write to a file).

## Concurrency Model

Each source and processing stage are executed via goroutines to allow messages to be streamed through the system (if all stages were executed on the main goroutine, each channel would have to drain before the next stage could begin). Channels are used to sync state between each goroutine, since send/recieve block until the other side is ready. The destination stream runs on the main goroutine to ensure the program doesn't terminate until that message stream finishes executing.

Source messages are additionally streamed via distinct goroutines to allow for parallel streaming.

To ensure graceful termination (exiting the program once all source channels are closed), we need to ensure the merged message stream is closed once all the source channels are closed. This is achieved using a `WaitGroup` (effectively Go's version of a semaphore): for each source stream, the semaphore is incremented and a goroutine is invoked to validate and send the message to `out`. Once the source channel is closed, the message range loop will terminate, and the semaphore will be decremented (via `WaitGroup.Done()`). Another goroutine blocks until the sempahore has been decremented to 0, and then closes the merged channel.

## Testing/Validation

## Limitations

    - Merge buffer overflow
    - client responsiblity (passing message through, error handling, closing channels)
    - messages must be strings.

## Future Improvements

    - Processing interface includes schema types/versions supported by that stage. System is responsible for validating messages into each stage, and passing them directly to the out stream so clients don't have to own this logic, and also allows for stronger validation between processing stages. This also means supported schema types/versions are dynamically set by clients rather than hardcoded in `Pipeline`.
