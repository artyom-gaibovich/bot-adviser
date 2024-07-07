package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Processor interface {
	Process(e Event) error
}

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
