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
	d := NewD("")

	kvput := d.DeclareChannel("kvput", KVPut{})
	kvputr := d.DeclareChannel("kvputr", KVPutResponse{})
	kvget := d.DeclareChannel("kvget", KVGet{})
	kvgetr := d.DeclareChannel("kvgetr", KVGetResponse{})

	if kvput == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvputr == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvget == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvgetr == nil {
		t.Errorf("expected non-nil channel")
	}

	kvstore := d.DeclareLMap("kvstore")

	kvstore.JoinMerge(kvput,
		func(k *KVPut) (interface{}, Lattice) { return k.Key, k.Val })

	kvputr.JoinMergeAsync(kvput,
		func(k *KVPut) *KVPutResponse {
			return &KVPutResponse{k.ReqId, k.ClientAddr, d.Addr}
		})

	kvgetr.JoinMergeAsync(kvget,
		func(k *KVGet) *KVGetResponse {
			return &KVGetResponse{k.ReqId, k.ClientAddr, d.Addr, k.Key,
				kvstore.At(k.Key)}
		})

	fmt.Printf("%#v\n", d)
}
