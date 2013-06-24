package gdec

import (
	"fmt"
	"reflect"
)

type D struct {
	Addr      string
	Relations map[string]Relation
	Joins     []*JoinDeclaration
	ticks     int64
}

type Relation interface {
	TupleType() reflect.Type
}

func NewD(addr string) *D {
	return &D{
		Addr:      addr,
		Relations: make(map[string]Relation),
		Joins:     []*JoinDeclaration{},
	}
}

func (d *D) DeclareRelation(name string, x Relation) Relation {
	if d.Relations[name] != nil {
		panic(fmt.Sprintf("relation redeclared, name: %s"+
			", relation: %#v", name, x))
	}
	d.Relations[name] = x
	return x
}

func (d *D) Join(vars ...interface{}) *JoinDeclaration {
	var r *Relation
	rt := reflect.TypeOf(r).Elem()

	var joinNum int
	var selectWhereFunc interface{}

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
			selectWhereFunc = x
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

	if selectWhereFunc != nil {
		mft := reflect.TypeOf(selectWhereFunc)
		if mft.NumOut() != 1 {
			panic(fmt.Sprintf("selectWhereFunc should return 1 value"+
				", selectWhereFunc: %v", mft))
		}
		if mft.NumIn() != joinNum {
			panic(fmt.Sprintf("selectWhereFunc should take %v args"+
				", selectWhereFunc: %v", joinNum, mft))
		}
		for i, x := range sources {
			rt := reflect.PtrTo(x.TupleType())
			if rt != mft.In(i) {
				panic(fmt.Sprintf("selectWhereFunc param #%v type"+
					"%v does not match, expected: %v, selectWhereFunc: %v",
					i, mft.In(i), rt, mft))
			}
		}
	}

	return &JoinDeclaration{
		d:               d,
		sources:         sources,
		selectWhereFunc: selectWhereFunc,
	}
}

func (d *D) JoinFlat(vars ...interface{}) *JoinDeclaration {
	jd := d.Join(vars...)
	jd.selectWhereFlat = true
	return jd
}

type JoinDeclaration struct {
	d               *D
	sources         []Relation
	selectWhereFunc interface{}
	selectWhereFlat bool
	async           bool
	into            Relation
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

	var out reflect.Type
	if jd.selectWhereFunc != nil {
		out = reflect.TypeOf(jd.selectWhereFunc).Out(0)
	} else if len(jd.sources) == 1 {
		out = reflect.PtrTo(jd.sources[0].TupleType())
	} else {
		panic(fmt.Sprintf("unexpected Into() join declaration: %#v", jd))
	}
	if jd.selectWhereFlat {
		if out != dt {
			panic(fmt.Sprintf("Into() param: %#v, type: %v, does not match"+
				" output type: %v", dest, dt, out))
		}
	} else {
		if out != jd.into.TupleType() &&
			out != reflect.PtrTo(jd.into.TupleType()) {
			panic(fmt.Sprintf("Into() param: %#v, type: %v, does not match"+
				" tuple type: %v", dest, dt, out))
		}
	}

	jd.d.Joins = append(jd.d.Joins, jd)
}
