package quic

import (
	"math"
	"testing"
)

func TestUint64ToInt32(t *testing.T) {
	if v, ok := uint64ToInt32(uint64(math.MaxInt32)); !ok || v != math.MaxInt32 {
		t.Fatalf("expected MaxInt32 conversion to succeed, got v=%d ok=%v", v, ok)
	}
	if _, ok := uint64ToInt32(uint64(math.MaxInt32) + 1); ok {
		t.Fatal("expected overflow conversion to fail")
	}
}
