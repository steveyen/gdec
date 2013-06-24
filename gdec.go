package gdec

import (
	"fmt"
	"reflect"
)

type D struct {
	Addr      string
	Relations map[string]Relation
	ticks     int64
}

type Joinable interface {
	TupleType() reflect.Type
}

type Channel struct {
	d *D
	t reflect.Type
}

type Relation interface{}

func NewD(addr string) *D {
	return &D{
		Addr:      addr,
		Relations: make(map[string]Relation),
	}
}

func (d *D) DeclareRelation(name string, x Relation) Relation {
	d.Relations[name] = x
	return x
}

func (d *D) Join(vars ...interface{}) *JoinDeclaration {
	var j *Joinable

	jt := reflect.TypeOf(j).Elem()

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
		} else if xt.Implements(jt) {
			joinNum = i + 1
		} else {
			panic(fmt.Sprintf("unexpected Join() param type: %#v, %v",
				x, xt))
		}
	}

	if mapFunc != nil {
		mft := reflect.TypeOf(mapFunc)
		if mft.NumOut() != 1 {
			panic(fmt.Sprintf("mapFunc should return 1 value, mapFunc: %v",
				mft))
		}
		if mft.NumIn() != joinNum {
			panic(fmt.Sprintf("mapFunc should take %v args, mapFunc: %v",
				joinNum, mft))
		}
		for i, x := range vars[0:joinNum] {
			j := x.(Joinable)
			jt := reflect.PtrTo(j.TupleType())
			if jt != mft.In(i) {
				panic(fmt.Sprintf("mapFunc param #%v type %v does not match, " +
					"expected: %v, mapFunc: %v", i, mft.In(i), jt, mft))
			}
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

func (d *D) DeclareChannel(name string, x interface{}) *Channel {
	return d.DeclareRelation(name, d.NewChannel(x)).(*Channel)

}

func (d *D) NewChannel(x interface{}) *Channel {
	return &Channel{d: d, t: reflect.TypeOf(x)}
}

func (c *Channel) TupleType() reflect.Type {
	return c.t
}
