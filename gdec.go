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
	for i, x := range vars {
		xt := reflect.TypeOf(x)
		switch xt.Kind() {
		case reflect.Func:
			if i < len(vars) - 1 {
				panic(fmt.Sprintf("func not last Join() param: %#v",
					vars))
			}
		case reflect.Ptr:
			switch xt.Elem().Kind() {
			case reflect.Interface:
			case reflect.Struct:
			default:
				panic(fmt.Sprintf("unexpected Join() param type: %#v",
					x))
			}
		default:
			panic(fmt.Sprintf("unexpected Join() param type: %#v",
				x))
		}
	}
	return nil
}

type JoinDeclaration struct {
	d *D
}

func (r *JoinDeclaration) Into(dest interface{}) {
}

func (r *JoinDeclaration) IntoAsync(dest interface{}) {
}
