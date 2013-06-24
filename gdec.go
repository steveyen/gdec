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

type Relation interface {
	TupleType() reflect.Type
}

func NewD(addr string) *D {
	return &D{
		Addr:      addr,
		Relations: make(map[string]Relation),
	}
}

func (d *D) DeclareRelation(name string, x Relation) Relation {
	if d.Relations[name] != nil {
		panic(fmt.Sprintf("relation redeclared, name: %s, relation: %#v", name, x))
	}
	d.Relations[name] = x
	return x
}

func (d *D) Join(vars ...interface{}) *JoinDeclaration {
	var r *Relation
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
		} else if xt.Implements(rt) {
			joinNum = i + 1
		} else {
			panic(fmt.Sprintf("unexpected Join() param type: %#v, %v",
				x, xt))
		}
	}

	sources := make([]Relation, joinNum)
	for i := 0; i < joinNum; i++ {
		sources[i] = vars[i].(Relation)
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
		for i, x := range sources {
			rt := reflect.PtrTo(x.TupleType())
			if rt != mft.In(i) {
				panic(fmt.Sprintf("mapFunc param #%v type %v does not match, "+
					"expected: %v, mapFunc: %v", i, mft.In(i), rt, mft))
			}
		}
	}

	return &JoinDeclaration{
		d:       d,
		sources: sources,
		mapFunc: mapFunc,
	}
}

func (d *D) JoinFlat(vars ...interface{}) *JoinDeclaration {
	jd := d.Join(vars...)
	jd.mapFlat = true
	return jd
}

type JoinDeclaration struct {
	d       *D
	sources []Relation
	mapFunc interface{}
	mapFlat bool
	async   bool
	into    Relation
}

func (jd *JoinDeclaration) IntoAsync(dest interface{}) {
	jd.async = true
	jd.Into(dest)
}

func (jd *JoinDeclaration) Into(dest interface{}) {
	var r *Relation
	rt := reflect.TypeOf(r).Elem()

	dt := reflect.TypeOf(dest)
	if !dt.Implements(rt) {
		panic(fmt.Sprintf("Into() param: %#v, type: %v, does not implement Relation",
			dest, dt))
	}

	jd.into = dest.(Relation)
}
