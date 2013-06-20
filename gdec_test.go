package gdec

import (
	"testing"
)

func TestNewD(t *testing.T) {
	if NewD() == nil {
		t.Errorf("expected D")
	}
}
