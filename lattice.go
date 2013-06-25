package gdec

import (
	"encoding/json"
	"reflect"
)

type Lattice interface{}

type LMap struct {
	d       *D
	scratch bool
}

type LMapEntry struct {
	Key string
	Val Lattice
}

type LSet struct {
	d       *D
	t       reflect.Type
	m       map[string]interface{}
	scratch bool
}

type LMax struct {
	d       *D
	v       int
	scratch bool
}

type LBool struct {
	d       *D
	v       bool
	scratch bool
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

func (d *D) NewLSet(x interface{}) *LSet {
	return &LSet{d: d, t: reflect.TypeOf(x), m: map[string]interface{}{}}
}

func (d *D) NewLMax() *LMax { return &LMax{d: d} }

func (d *D) NewLBool() *LBool { return &LBool{d: d} }

func (m *LMap) TupleType() reflect.Type {
	var x *LMapEntry
	return reflect.TypeOf(x).Elem()
}

func (m *LSet) TupleType() reflect.Type {
	return m.t
}

func (m *LMax) TupleType() reflect.Type {
	return reflect.TypeOf(0)
}

func (m *LBool) TupleType() reflect.Type {
	var x bool
	return reflect.TypeOf(x)
}

func (m *LMap) DeclareScratch() {
	m.scratch = true
}

func (m *LSet) DeclareScratch() {
	m.scratch = true
}

func (m *LMax) DeclareScratch() {
	m.scratch = true
}

func (m *LBool) DeclareScratch() {
	m.scratch = true
}

func (m *LMap) startTick() {
	if m.scratch {
		// TODO.
	}
}

func (m *LSet) startTick() {
	if m.scratch {
		m.m = map[string]interface{}{}
	}
}

func (m *LMax) startTick() {
	if m.scratch {
		m.v = 0
	}
}

func (m *LBool) startTick() {
	if m.scratch {
		m.v = false
	}
}

func (m *LMap) Add(v interface{}) bool {
	panic("LMap.Add unimplemented")
	return false
}

func (m *LSet) Add(v interface{}) bool {
	if v == nil {
		panic("unexpected nil during LSet.Add")
	}
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	if string(j) == "null" {
		panic("unexpected null during LSet.Add")
	}
	js := string(j)
	_, changed := m.m[js]
	m.m[js] = v
	return changed
}

func (m *LMax) Add(v interface{}) bool {
	vi := v.(int)
	if m.v < vi {
		m.v = vi
		return true
	}
	return false
}

func (m *LBool) Add(v interface{}) bool {
	old := m.v
	m.v = m.v || v.(bool)
	return m.v == old
}

func (m *LMap) Merge(rel Relation) bool {
	panic("LMap.Merge unimplemented")
	return false
}

func (m *LSet) Merge(rel Relation) bool {
	changed := false
	r := rel.(*LSet)
	for _, v := range r.m {
		changed = changed || m.Add(v)
	}
	return changed
}

func (m *LMax) Merge(rel Relation) bool {
	return m.Add(rel.(*LMax).v)
}

func (m *LBool) Merge(rel Relation) bool {
	return m.Add(rel.(*LBool).v)
}

func (m *LMap) Scan() chan interface{} {
	panic("LMap.Scan unimplemented")
	return nil
}

func (m *LSet) Scan() chan interface{} {
	ch := make(chan interface{})
	go func() {
		for _, v := range m.m {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func (m *LMax) Scan() chan interface{} {
	ch := make(chan interface{})
	go func() {
		ch <- m.v
		close(ch)
	}()
	return ch
}

func (m *LBool) Scan() chan interface{} {
	ch := make(chan interface{})
	go func() {
		ch <- m.v
		close(ch)
	}()
	return ch
}

func (m *LMap) At(key string) Lattice {
	return nil
}

func (m *LMap) Snapshot() *LMap {
	return nil
}

func (m *LSet) Size() int {
	return len(m.m)
}

func (m *LMax) Int() int {
	return m.v
}

func (m *LBool) Bool() bool {
	return m.v
}
