package gdec

import (
	"reflect"
)

type Lattice interface{}

type Channel struct {
	t reflect.Type
}

type D struct {
	channels map[string]*Channel
}

func (d *D) RegisterChannel(name string, x interface{}) *Channel {
	c := NewChannel(x)
	d.channels[name] = c
	return c
}
func NewChannel(x interface{}) *Channel {
	return &Channel{t: reflect.TypeOf(x)}
}

func NewD() *D {
	return &D{
		channels: make(map[string]*Channel),
	}
}
