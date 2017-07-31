package main

import (
	"testing"
)

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	str := randStringBytesMaskImprSrc(64)
	if len(str) != 64 {
		t.Fatalf("randStringBytesMaskImprSrc(%d) == %d, want %d", 64, len(str), 64)
	}
}
