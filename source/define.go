package source

import "context"

type SourceInstance interface {
	GetName() string
	Start(context.Context, chan<- []byte)
	Alive() bool
	Stop()
}

type NewSourceInstanceFunc func(interface{}) (SourceInstance, error)

var (
	instanceFactory = map[string]NewSourceInstanceFunc{}
)

func RegisterSourceInstance(name string, f NewSourceInstanceFunc) {
	if f != nil && name != "" {
		instanceFactory[name] = f
	}
}
