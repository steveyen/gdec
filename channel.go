package gdec

import (
	"reflect"
)

type Channel struct {
	d       *D
	t       reflect.Type
	scratch bool
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

func (c *Channel) DeclareScratch() {
	c.scratch = true
}

func (c *Channel) startTick() {
	if c.scratch {
		// TODO.
	}
}

func (c *Channel) Add(v interface{}) bool {
	panic("Channel.Add() unimplemented")
	return false
}

func (c *Channel) Merge(rel Relation) bool {
	panic("Channel.Merge unimplemented")
	return false
}

func (c *Channel) Scan() chan interface{} {
	panic("Channel.Scan() unimplemented")
}
