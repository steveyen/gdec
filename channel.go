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
