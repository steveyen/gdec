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

func TestQuorum(t *testing.T) {
	d := QuorumInit(NewD(""), "", 5, "")
	fmt.Printf("%#v\n", d)
}

func TestShortestPath(t *testing.T) {
	d := ShortestPathInit(NewD(""), "")
	links := d.Relations["ShortestPathLink"].(*LSet)
	paths := d.Relations["ShortestPath"].(*LSet)

	links.Add(&ShortestPathLink{From: "a", To: "b", Cost: 1})
	links.Add(&ShortestPathLink{From: "b", To: "c", Cost: 1})
	if links.Size() != 2 {
		t.Errorf("expected 2 links, got: %v", links.Size())
	}
	if paths.Size() != 0 {
		t.Errorf("expected 0 links, got: %v", paths.Size())
	}

	d.Tick()
}
