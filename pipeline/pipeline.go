package pipeline

type Pipeline struct {
	Sources []Source
	Processings []Processing
	// Single destination is in scope.
	Destination Destination
}

func (p *Pipeline) Run() {
	c := make(chan Message)
	// read-one/write-one line-by-line flow.

	// does everythign need to have multiple channels?
	for _, source := range p.Sources {
		go source.Messages(c)
	}
	p.Destination.Messages(c)
}