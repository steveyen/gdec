package gdec

func (d *D) DeclareChannel(name string, x interface{}) *LSet {
	c := d.DeclareLSet(name, x)
	c.channel = true
	return c
}
