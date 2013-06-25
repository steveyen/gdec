package gdec

import (
	"reflect"
)

type Channel struct {
	d *D
	t reflect.Type
}

func (d *D) DeclareChannel(name string, x interface{}) *Channel {
	return d.DeclareRelation(name, d.NewChannel(x)).(*Channel)
}

func (d *D) NewChannel(x interface{}) *Channel {
	return &Channel{d: d, t: reflect.TypeOf(x)}
}

func (c *Channel) TupleType() reflect.Type {
	return c.t
}

func (c *Channel) Add(v interface{}) {
	panic("Channel.Add() unimplemented")
}

func (c *Channel) Merge(rel Relation) {
	panic("Channel.Merge unimplemented")
}

func (c *Channel) Scan() chan interface{} {
	panic("Channel.Scan() unimplemented")
}
