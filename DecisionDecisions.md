## Overall architecture



## Assumptions
- Synthesize multiple sources; only support one destination.

## Error Handling

Errors may occur when streaming data in the message source, or due to unexpected data in the processing stages.

To keep the message stream simple, each stage is currently responsible for handling any errors it incurs.
