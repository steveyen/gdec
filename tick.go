package gdec

import (
	"fmt"
	"reflect"
)

func (d *D) Tick() {
	d.tickBefore()
	d.tickCore()
	d.ticks++
}

func (d *D) tickBefore() {
	// TODO: Incorporate periodics.
	// TODO: Incorporate network.
	// Incorporate next.
	for name, arr := range d.next {
		r := d.Relations[name]
		if r == nil {
			panic(fmt.Sprintf("unknown relation: %s", name))
		}
		for _, x := range arr {
			r.Add(x)
		}
	}
	d.next = make(map[string][]interface{})
}

func (d *D) tickCore() {
	for _, jd := range d.Joins {
		jd.executeJoinInto()
	}
}

func (jd *JoinDeclaration) executeJoinInto() {
	numSources := len(jd.sources)

	join := make([]interface{}, numSources)

	type joinAccum struct {
		into Relation
		val interface{}
		add bool
	}
	accums := []joinAccum{}

	accum := func() {
		if jd.selectWhereFunc != nil {
			values := make([]reflect.Value, numSources)
			for i, x := range join {
				values[i] = reflect.ValueOf(x)
			}
			ft := reflect.ValueOf(jd.selectWhereFunc)
			out := ft.Call(values)
			if out == nil || len(out) != 1 {
				panic(fmt.Sprintf("unexpected # out results: %#v", out))
			}
			if !out[0].IsNil() {
				out0 := out[0].Interface()
				if out0 != nil {
					if jd.selectWhereFlat {
						accums = append(accums, joinAccum{jd.into, out0, false})
					} else {
						accums = append(accums, joinAccum{jd.into, out0, true})
					}
				}
			}
		} else if len(join) == 1 {
			accums = append(accums, joinAccum{jd.into, join[0], true})
		} else {
			panic("could not send join output into receiver")
		}
	}

	var joiner func(int)
	joiner = func(pos int) {
		if pos < numSources {
			for tuple := range(jd.sources[pos].Scan()) {
				if tuple == nil {
					panic("Scan() gave nil tuple")
				}
				join[pos] = tuple
				joiner(pos + 1)
			}
		} else {
			accum()
		}
	}
	joiner(0)

	for _, a := range accums {
		if a.add {
			a.into.Add(a.val)
		} else {
			a.into.Merge(a.val.(Relation))
		}
	}
}
