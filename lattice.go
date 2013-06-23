package gdec

type Lattice interface{}

type LMap struct {
	d *D
}

type LSet struct {
	d *D
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

func (d *D) DeclareLSet(name string, t interface{}) *LSet {
	return d.DeclareRelation(name, d.NewLSet(t)).(*LSet)
}

func (d *D) DeclareLMax(name string) *LMax {
	return d.DeclareRelation(name, d.NewLMax()).(*LMax)
}

func (d *D) DeclareLBool(name string) *LBool {
	return d.DeclareRelation(name, d.NewLBool()).(*LBool)
}

func (d *D) NewLMap() *LMap { return &LMap{d: d} }

func (d *D) NewLSet(t interface{}) *LSet { return &LSet{d: d} }

func (d *D) NewLMax() *LMax { return &LMax{d: d} }

func (d *D) NewLBool() *LBool { return &LBool{d: d} }

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
