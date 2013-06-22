package gdec

import (
	"reflect"
)

type D struct {
	Addr     string
	Channels map[string]*Channel
	Lattices map[string]Lattice
}

type Channel struct {
	d *D
	t reflect.Type
}

type JoinResult struct {
	d *D
}

func NewD(addr string) *D {
	return &D{
		Addr:     addr,
		Channels: make(map[string]*Channel),
		Lattices: make(map[string]Lattice),
	}
}

func (d *D) DeclareChannel(name string, x interface{}) *Channel {
	c := d.NewChannel(x)
	d.Channels[name] = c
	return c
}

func (d *D) NewChannel(x interface{}) *Channel {
	return &Channel{d: d, t: reflect.TypeOf(x)}
}

func (d *D) Join(v ...interface{}) *JoinResult {
	return nil
}

func (r *JoinResult) Into(dest interface{}) {
}

func (r *JoinResult) IntoAsync(dest interface{}) {
}
