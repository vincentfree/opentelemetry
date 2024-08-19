package providerconfig

import "slices"

type Execution uint8

//sync async

const (
	Sync Execution = iota + 1
	Async
)

var (
	executionTypes = []Execution{Sync, Async}
)

func (e Execution) String() string {
	switch e {
	case Sync:
		return "Sync"
	case Async:
		return "Async"
	default:
		return "undefined"
	}
}

func (e Execution) IsValid() bool {
	return slices.Contains(executionTypes, e)
}
