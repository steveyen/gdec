package gdec

import (
	"reflect"
)

type Lattice interface{}

type D struct {
	Addr string
	Channels map[string]*Channel
	Lattices map[string]Lattice
}

func NewD(addr string) *D {
	return &D{
		Addr: addr,
		Channels: make(map[string]*Channel),
		Lattices: make(map[string]Lattice),
	}
}

func (d *D) DeclareChannel(name string, x interface{}) *Channel {
	c := d.NewChannel(x)
	d.Channels[name] = c
	return c
}

type Channel struct {
	d *D
	t reflect.Type
}

func (d *D) NewChannel(x interface{}) *Channel {
	return &Channel{d: d, t: reflect.TypeOf(x)}
}

func (c *Channel) JoinMergeAsync(v ...interface{}) {
}

func (d *D) DeclareLMap(name string) *LMap {
	m := d.NewLMap()
	d.Lattices[name] = m
	return m
}

type LMap struct {
	d *D
}

func (d *D) NewLMap() *LMap {
	return &LMap{d: d}
}

func (m *LMap) At(key string) Lattice {
	return nil
}

func (m *LMap) JoinMerge(v ...interface{}) {
}
