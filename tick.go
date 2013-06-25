package gdec

import (
	"fmt"
	"reflect"
)

type relationChange struct {
	into Relation
	arg  interface{} // Arg for Add/Merge() call.
	add  bool        // Use Add() versus Merge().
}

func (d *D) Tick() {
	d.tickBefore()
	d.tickCore()
	d.ticks++
}

func (d *D) tickBefore() {
	// TODO: Incorporate periodics.
	// TODO: Incorporate network.

	applyRelationChanges(d.next)
	d.next = d.next[0:0]
}

func (d *D) tickCore() {
	for _, jd := range d.Joins {
		jd.executeJoinInto()
	}
}

func (jd *JoinDeclaration) executeJoinInto() {
	numSources := len(jd.sources)

	join := make([]interface{}, numSources)

	accums := []relationChange{}

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
						accums = append(accums,
							relationChange{jd.into, out0, false})
					} else {
						accums = append(accums,
							relationChange{jd.into, out0, true})
					}
				}
			}
		} else if len(join) == 1 {
			if join[0] != nil {
				accums = append(accums,
					relationChange{jd.into, join[0], true})
			}
		} else {
			panic("could not send join output into receiver")
		}
	}

	var joiner func(int)
	joiner = func(pos int) {
		if pos < numSources {
			for tuple := range jd.sources[pos].Scan() {
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

	applyRelationChanges(accums)
}

func applyRelationChanges(changes []relationChange) {
	for _, c := range changes {
		if c.add {
			c.into.Add(c.arg)
		} else {
			c.into.Merge(c.arg.(Relation))
		}
	}
}
