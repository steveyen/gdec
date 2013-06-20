package gdec

import (
	"fmt"
	"testing"
)

func TestNewD(t *testing.T) {
	if NewD() == nil {
		t.Errorf("expected D")
	}
}

func TestKV(t *testing.T) {
	d := NewD()

	kvputs := d.RegisterChannel("kvputs", KVPut{})
	kvputr := d.RegisterChannel("kvput_responses", KVPutResponse{})
	kvgets := d.RegisterChannel("kvgets", KVGet{})
	kvgetr := d.RegisterChannel("kvget_responses", KVGetResponse{})

	if kvputs == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvputr == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvgets == nil {
		t.Errorf("expected non-nil channel")
	}
	if kvgetr == nil {
		t.Errorf("expected non-nil channel")
	}

	fmt.Printf("%#v\n", d)
}
