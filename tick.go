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
	for _, r := range d.Relations {
		r.startTick()
	}

	// TODO: Incorporate periodics.
	// TODO: Incorporate network.

	applyRelationChanges(d.next) // Apply pending data from last tick.
	d.next = d.next[0:0]

	d.tickMain()
	d.ticks++

	// TODO: Emit to network.
}

func (d *D) tickMain() {
	for { // TODO: Hugely naive, inefficient, simple implementation.
		changed := false
		for _, jd := range d.Joins {
			d.next, d.immediate = jd.executeJoinInto(d.next, d.immediate)
			changed = changed || applyRelationChanges(d.immediate)
			d.immediate = d.immediate[0:0]
		}
		if !changed {
			return
		}
	}
}

func (jd *joinDeclaration) executeJoinInto(next, immediate []relationChange) (
	nextOut, immediateOut []relationChange) {
	numSources := len(jd.sources)

	join := make([]interface{}, numSources)
	values := make([]reflect.Value, numSources)

	selectWhere := func() *relationChange {
		if jd.selectWhereFunc != nil {
			for i, x := range join {
				values[i] = reflect.ValueOf(x)
			}
			ft := reflect.ValueOf(jd.selectWhereFunc)
			out := ft.Call(values)
			if out == nil || len(out) != 1 {
				panic(fmt.Sprintf("unexpected # out results: %#v", out))
			}
			if out[0].IsValid() && !isNil(out[0]) {
				out0 := out[0].Interface()
				if out0 != nil {
					if jd.selectWhereFlat {
						return &relationChange{jd.into, out0, false}
					} else {
						return &relationChange{jd.into, out0, true}
					}
				}
			}
		} else if len(join) == 1 {
			if join[0] != nil {
				return &relationChange{jd.into, join[0], true}
			}
		} else {
			panic("could not send join output into receiver")
		}
		return nil
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
			res := selectWhere()
			if res != nil {
				if jd.async {
					next = append(next, *res)
				} else {
					immediate = append(immediate, *res)
				}
			}
		}
	}
	joiner(0)

	return next, immediate
}

func applyRelationChanges(changes []relationChange) bool {
	changed := false
	for _, c := range changes {
		if c.add {
			changed = changed || c.into.DirectAdd(c.arg)
		} else {
			changed = changed || c.into.DirectMerge(c.arg.(Relation))
		}
	}
	return changed
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}
