package gdec

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Lattice interface {
	DirectMerge(rel Relation) bool
	Snapshot() Lattice
}

type LMap struct {
	name    string
	d       *D
	m       map[string]Lattice
	scratch bool
}

type LMapEntry struct {
	Key string
	Val Lattice
}

type LSet struct {
	name    string
	d       *D
	t       reflect.Type
	m       map[string]interface{}
	scratch bool
	channel bool // When true, this LSet was declared as a channel.
}

type LMax struct {
	name    string
	d       *D
	v       int
	scratch bool
}

type LMaxString struct {
	name    string
	d       *D
	v       string
	scratch bool
}

type LBool struct {
	name    string
	d       *D
	v       bool
	scratch bool
}

func (d *D) DeclareLMap(name string) *LMap {
	m := d.NewLMap()
	m.name = name
	return d.DeclareRelation(name, m).(*LMap)
}

func (d *D) DeclareLSet(name string, x interface{}) *LSet {
	m := d.NewLSet(reflect.TypeOf(x))
	m.name = name
	return d.DeclareRelation(name, m).(*LSet)
}

func (d *D) DeclareLMax(name string) *LMax {
	m := d.NewLMax()
	m.name = name
	return d.DeclareRelation(name, m).(*LMax)
}

func (d *D) DeclareLMaxString(name string) *LMaxString {
	m := d.NewLMaxString()
	m.name = name
	return d.DeclareRelation(name, m).(*LMaxString)
}

func (d *D) DeclareLBool(name string) *LBool {
	m := d.NewLBool()
	m.name = name
	return d.DeclareRelation(name, m).(*LBool)
}

func (d *D) NewLMap() *LMap { return &LMap{d: d, m: map[string]Lattice{}} }

func (d *D) NewLSet(t reflect.Type) *LSet {
	return &LSet{d: d, t: t, m: map[string]interface{}{}}
}

func (d *D) NewLMax() *LMax { return &LMax{d: d} }

func (d *D) NewLMaxString() *LMaxString { return &LMaxString{d: d} }

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

func (m *LMaxString) TupleType() reflect.Type {
	return reflect.TypeOf("")
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

func (m *LMaxString) DeclareScratch() {
	m.scratch = true
}

func (m *LBool) DeclareScratch() {
	m.scratch = true
}

func (m *LMap) startTick() {
	if m.scratch {
		m.m = map[string]Lattice{}
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

func (m *LMaxString) startTick() {
	if m.scratch {
		m.v = ""
	}
}

func (m *LBool) startTick() {
	if m.scratch {
		m.v = false
	}
}

func (m *LMap) DirectAdd(v interface{}) bool {
	if v == nil {
		panic("unexpected nil during LMap.DirectAdd")
	}
	e := v.(*LMapEntry)
	o, _ := m.m[e.Key]
	if o != nil {
		changed := o.DirectMerge(e.Val.(Relation))
		m.m[e.Key] = o
		return changed
	}
	m.m[e.Key] = e.Val
	return true
}

func (m *LSet) DirectAdd(v interface{}) bool {
	if v == nil {
		panic("unexpected nil during LSet.DirectAdd")
	}
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	if string(j) == "null" {
		panic(fmt.Sprintf("unexpected null during LSet.DirectAdd"+
			", v: %#v, LSet.name: %s", v, m.name))
	}
	js := string(j)
	_, exists := m.m[js]
	m.m[js] = v
	return !exists
}

func (m *LMax) DirectAdd(v interface{}) bool {
	vi := v.(int)
	if m.v < vi {
		m.v = vi
		return true
	}
	return false
}

func (m *LMaxString) DirectAdd(v interface{}) bool {
	vs := v.(string)
	if m.v < vs {
		m.v = vs
		return true
	}
	return false
}

func (m *LBool) DirectAdd(v interface{}) bool {
	old := m.v
	m.v = m.v || v.(bool)
	return m.v != old
}

func (m *LMap) DirectMerge(rel Relation) bool {
	changed := false
	r := rel.(*LMap)
	for k, v := range r.m {
		changed = m.DirectAdd(&LMapEntry{k, v}) || changed
	}
	return changed
}

func (m *LSet) DirectMerge(rel Relation) bool {
	changed := false
	r := rel.(*LSet)
	for _, v := range r.m {
		changed = m.DirectAdd(v) || changed
	}
	return changed
}

func (m *LMax) DirectMerge(rel Relation) bool {
	return m.DirectAdd(rel.(*LMax).v)
}

func (m *LMaxString) DirectMerge(rel Relation) bool {
	return m.DirectAdd(rel.(*LMaxString).v)
}

func (m *LBool) DirectMerge(rel Relation) bool {
	return m.DirectAdd(rel.(*LBool).v)
}

func (m *LMap) Scan() chan interface{} {
	ch := make(chan interface{})
	go func() {
		for k, v := range m.m {
			ch <- &LMapEntry{k, v}
		}
		close(ch)
	}()
	return ch
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

func (m *LMaxString) Scan() chan interface{} {
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

func (m *LMap) Snapshot() Lattice {
	s := m.d.NewLMap()
	for k, v := range m.m {
		s.m[k] = v.Snapshot()
	}
	return s
}

func (m *LSet) Snapshot() Lattice {
	s := m.d.NewLSet(m.t)
	for k, v := range m.m {
		switch v.(type) {
		case Lattice:
			s.m[k] = v.(Lattice).Snapshot()
		default:
			s.m[k] = v
		}
	}
	return s
}

func (m *LMax) Snapshot() Lattice {
	s := m.d.NewLMax()
	s.v = m.v
	return s
}

func (m *LMaxString) Snapshot() Lattice {
	s := m.d.NewLMaxString()
	s.v = m.v
	return s
}

func (m *LBool) Snapshot() Lattice {
	s := m.d.NewLBool()
	s.v = m.v
	return s
}

func (m *LMap) At(key string) Lattice {
	v, _ := m.m[key]
	return v
}

func (m *LSet) Contains(v interface{}) bool {
	if v == nil {
		panic("unexpected nil during LSet.Contains")
	}
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	if string(j) == "null" {
		panic("unexpected null during LSet.Contains")
	}
	js := string(j)
	_, ok := m.m[js]
	return ok
}

func (m *LSet) Size() int {
	return len(m.m)
}

func (m *LMax) Int() int {
	return m.v
}

func (m *LMaxString) String() string {
	return m.v
}

func (m *LBool) Bool() bool {
	return m.v
}

func NewLSetOne(d *D, v interface{}) *LSet { // Helper creator for a 1 item LSet.
	s := d.NewLSet(reflect.TypeOf(v))
	s.DirectAdd(v)
	return s
}

func NewLBool(d *D, v bool) *LBool { // Helper creator for an initialized LBool.
	s := d.NewLBool()
	s.DirectAdd(v)
	return s
}
