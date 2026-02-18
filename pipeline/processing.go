package pipeline
// Validation and transformation for now - may be seperate interfaces later
type Processing interface {
	// channels from all sources need to be available to the processesing stage.
	// Or actually from previous processing stages as well.
	// Messages themselves likely need an indentifer, and then sources need to
	// describe themselves and their schema.
	// New schemas may be defined by processing stages as they provide transformations to later stages.
	Process(chan Message)
}

// some processing stages transform (modify, reduce filter) messages based on a single schema.
// others may combine multiple messages (i.e., with the same timestamp and site internal id) into one representation. 


