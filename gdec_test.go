package gdec

import (
	"fmt"
	"testing"
)

func TestNewD(t *testing.T) {
	if NewD("") == nil {
		t.Errorf("expected D")
	}
}

func TestKV(t *testing.T) {
	d := KVInit(NewD(""), "")
	fmt.Printf("%#v\n", d)
}

func TestReplicatedKV(t *testing.T) {
	d := ReplicatedKVInit(NewD(""), "")
	fmt.Printf("%#v\n", d)
}

func TestTally(t *testing.T) {
	d := TallyInit(NewD("tallyTest"), "")

	tvote := d.Relations["TallyVote"].(*LSet)
	tneed := d.Relations["TallyNeed"].(*LMax)
	tdone := d.Relations["TallyDone"].(*LBool)

	if tdone.Bool() {
		t.Errorf("shouldn't have done tally already")
	}
	d.Tick()
	if !tdone.Bool() {
		t.Errorf("should have done 0 tally already")
	}

	if !tneed.DirectAdd(2) {
		t.Errorf("expected tneed to change")
	}
	if tneed.Int() != 2 {
		t.Errorf("expected tneed to be 2")
	}
	d.Tick()
	if tdone.Bool() {
		t.Errorf("should not have done 2 tally already")
	}

	d.AddNext(tvote, "a")
	d.Tick()
	if tdone.Bool() {
		t.Errorf("should not have done 2 tally already")
	}

	d.AddNext(tvote, "a")
	d.Tick()
	if tdone.Bool() {
		t.Errorf("should not have done 2 tally already")
	}

	d.AddNext(tvote, "b")
	d.Tick()
	if !tdone.Bool() {
		t.Errorf("should have done 2 tally already")
	}

	d.AddNext(tvote, "b")
	d.Tick()
	if !tdone.Bool() {
		t.Errorf("should have done 2 tally already")
	}

	d.AddNext(tvote, "c")
	d.Tick()
	if !tdone.Bool() {
		t.Errorf("should stay done at 2 tally already")
	}
}

func TestMultiTally(t *testing.T) {
	d := MultiTallyInit(NewD("multiTallyTest"), "")

	tvote := d.Relations["MultiTallyVote"].(*LSet)
	tneed := d.Relations["MultiTallyNeed"].(*LMax)
	tdone := d.Relations["MultiTallyDone"].(*LMap)

	if !tneed.DirectAdd(2) {
		t.Errorf("expected tneed to change")
	}
	if tneed.Int() != 2 {
		t.Errorf("expected tneed to be 2")
	}
	d.Tick()
	if tdone.At("A") != nil {
		t.Errorf("should not have done for A")
	}

	d.AddNext(tvote, &MultiTallyVote{"A", "a0"})
	d.Tick()
	if tdone.At("A").(*LBool).Bool() {
		t.Errorf("should not have done for A")
	}

	d.AddNext(tvote, &MultiTallyVote{"A", "a0"})
	d.Tick()
	if tdone.At("A").(*LBool).Bool() {
		t.Errorf("should not have done for A")
	}
	if tdone.At("B") != nil {
		t.Errorf("should not have done for B")
	}

	d.AddNext(tvote, &MultiTallyVote{"B", "b0"})
	d.AddNext(tvote, &MultiTallyVote{"A", "a1"})
	d.Tick()
	if !tdone.At("A").(*LBool).Bool() {
		t.Errorf("should be done for A")
	}
	if tdone.At("B").(*LBool).Bool() {
		t.Errorf("should not have done for B")
	}
}

func TestShortestPath(t *testing.T) {
	d := ShortestPathInit(NewD(""), "")
	links := d.Relations["ShortestPathLink"].(*LSet)
	paths := d.Relations["ShortestPath"].(*LSet)

	links.DirectAdd(&ShortestPathLink{From: "a", To: "b", Cost: 10})
	links.DirectAdd(&ShortestPathLink{From: "b", To: "c", Cost: 10})
	if links.Size() != 2 {
		t.Errorf("expected 2 links, got: %v", links.Size())
	}
	if paths.Size() != 0 {
		t.Errorf("expected 0 links, got: %v", paths.Size())
	}

	d.Tick()
	if d.ticks != 1 {
		t.Errorf("expected 1 ticks, got: %v", d.ticks)
	}
	if paths.Size() != 3 {
		t.Errorf("expected 3 links, got: %v, paths: %#v", paths.Size(), paths.m)
	}
	if !paths.Contains(&ShortestPath{From: "a", To: "c", Next: "b", Cost: 20}) {
		t.Errorf("expected paths to contain a->b")
	}

	d = ShortestPathInit(NewD(""), "")
	links = d.Relations["ShortestPathLink"].(*LSet)
	paths = d.Relations["ShortestPath"].(*LSet)
	links.DirectAdd(&ShortestPathLink{From: "a", To: "b", Cost: 10})
	links.DirectAdd(&ShortestPathLink{From: "b", To: "c", Cost: 10})
	links.DirectAdd(&ShortestPathLink{From: "a", To: "b", Cost: 1})
	d.Tick()
	if paths.Size() != 5 {
		t.Errorf("expected 5 links, got: %v, paths: %#v", paths.Size(), paths.m)
	}
	if !paths.Contains(&ShortestPath{From: "a", To: "c", Next: "b", Cost: 20}) {
		t.Errorf("expected paths to contain a->b")
	}
	if !paths.Contains(&ShortestPath{From: "a", To: "c", Next: "b", Cost: 11}) {
		t.Errorf("expected paths to contain a->b")
	}
	if paths.Contains(&ShortestPath{From: "a", To: "c", Next: "b", Cost: 1}) {
		t.Errorf("expected paths to to not contain a->b at the wrong cost")
	}
}
