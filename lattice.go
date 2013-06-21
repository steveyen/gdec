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

func (d *D) NewLMap() *LMap { return &LMap{d: d} }
func (d *D) NewLSet() *LSet { return &LSet{d: d} }
func (d *D) NewLMax() *LMax { return &LMax{d: d} }
func (d *D) NewLBool() *LBool { return &LBool{d: d} }

func (d *D) DeclareLMap(name string) *LMap {
	m := d.NewLMap()
	d.Lattices[name] = m
	return m
}

func (d *D) DeclareLSet(name string, x interface{}) *LSet {
	m := d.NewLSet()
	d.Lattices[name] = m
	return m
}

func (d *D) DeclareLMax(name string) *LMax {
	m := d.NewLMax()
	d.Lattices[name] = m
	return m
}

func (d *D) DeclareLBool(name string) *LBool {
	m := d.NewLBool()
	d.Lattices[name] = m
	return m
}

func (m *LMap) At(key string) Lattice {
	return nil
}

func (m *LMap) JoinUpdate(v ...interface{}) {
}

func (m *LSet) JoinUpdate(v ...interface{}) {
}

func (m *LSet) Size() int {
	return 0
}

func (m *LMax) Update(v ...interface{}) {
}

func (m *LMax) Int() int {
	return 0
}

func (m *LBool) Update(v ...interface{}) {
}

func (m *LBool) Bool() bool {
	return false
}
