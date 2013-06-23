package gdec

import (
	"reflect"
)

type D struct {
	Addr      string
	Channels  map[string]*Channel
	Relations map[string]Relation
}

type Channel struct {
	d *D
	t reflect.Type
}

type Relation interface{}

func NewD(addr string) *D {
	return &D{
		Addr:      addr,
		Channels:  make(map[string]*Channel),
		Relations: make(map[string]Relation),
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

func (d *D) DeclareRelation(name string, x Relation) Relation {
	d.Relations[name] = x
	return x
}

func (d *D) Join(v ...interface{}) *JoinRelation {
	return nil
}

type JoinRelation struct {
	d *D
}

func (r *JoinRelation) Into(dest interface{}) {
}

func (r *JoinRelation) IntoAsync(dest interface{}) {
}
