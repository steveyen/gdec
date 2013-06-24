package gdec

import (
	"fmt"
	"reflect"
)

type D struct {
	Addr      string
	Channels  map[string]*Channel
	Relations map[string]Relation
	ticks     int64
}

type Channel struct {
	d *D
	t reflect.Type
}

func (d *D) NewChannel(x interface{}) *Channel {
	return &Channel{d: d, t: reflect.TypeOf(x)}
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

func (d *D) DeclareRelation(name string, x Relation) Relation {
	d.Relations[name] = x
	return x
}

func (d *D) Join(vars ...interface{}) *JoinDeclaration {
	var c *Channel
	var r *Relation

	ct := reflect.TypeOf(c).Elem()
	rt := reflect.TypeOf(r).Elem()

	var joinNum int
	var mapFunc interface{}

	for i, x := range vars {
		if x == nil {
			panic("nil passed as Join() param")
		}
		xt := reflect.TypeOf(x)
		if xt.Kind() == reflect.Func {
			if i < len(vars)-1 {
				panic(fmt.Sprintf("func not last Join() param: %#v",
					vars))
			}
			mapFunc = x
		} else if xt.Kind() == reflect.Ptr {
			xt = xt.Elem()
			if !xt.AssignableTo(ct) && !xt.Implements(rt) {
				panic(fmt.Sprintf("unexpected Join() param pointer type: %#v",
					x))
			}
			joinNum = i + 1
		} else {
			panic(fmt.Sprintf("unexpected Join() param type: %#v, %v",
				x, xt))
		}
	}
	return &JoinDeclaration{
		d:       d,
		sources: vars[0:joinNum],
		mapFunc: mapFunc,
	}
}

type JoinDeclaration struct {
	d       *D
	sources []interface{}
	mapFunc interface{}
}

func (r *JoinDeclaration) Into(dest interface{}) {
}

func (r *JoinDeclaration) IntoAsync(dest interface{}) {
}
