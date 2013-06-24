package gdec

import (
	"reflect"
)

type Lattice interface{}

type LMap struct {
	d *D
}

type LSet struct {
	d *D
	t reflect.Type
}

type LMax struct {
	d *D
}

type LBool struct {
	d *D
}

func (d *D) DeclareLMap(name string) *LMap {
	return d.DeclareRelation(name, d.NewLMap()).(*LMap)
}

func (d *D) DeclareLSet(name string, x interface{}) *LSet {
	return d.DeclareRelation(name, d.NewLSet(x)).(*LSet)
}

func (d *D) DeclareLMax(name string) *LMax {
	return d.DeclareRelation(name, d.NewLMax()).(*LMax)
}

func (d *D) DeclareLBool(name string) *LBool {
	return d.DeclareRelation(name, d.NewLBool()).(*LBool)
}

func (d *D) NewLMap() *LMap { return &LMap{d: d} }

func (d *D) NewLSet(x interface{}) *LSet { return &LSet{d: d, t: reflect.TypeOf(x)} }

func (d *D) NewLMax() *LMax { return &LMax{d: d} }

func (d *D) NewLBool() *LBool { return &LBool{d: d} }

func (m *LMap) TupleType() reflect.Type {
	var x *Lattice
	return reflect.TypeOf(x).Elem()
}

func (m *LSet) TupleType() reflect.Type {
	return m.t
}

func (m *LMax) TupleType() reflect.Type {
	var x int
	return reflect.TypeOf(x)
}

func (m *LBool) TupleType() reflect.Type {
	var x bool
	return reflect.TypeOf(x)
}

func (m *LMap) At(key string) Lattice {
	return nil
}

func (m *LMap) Snapshot() *LMap {
	return nil
}

func (m *LSet) Size() int {
	return 0
}

func (m *LMax) Int() int {
	return 0
}

func (m *LBool) Bool() bool {
	return false
}
