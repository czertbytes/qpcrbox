package main

import (
	"testing"
)

func TestContentValidity(t *testing.T) {
	var md = AB7300{Content: "aaaa", Mock: "aaa"}

	if _, err := md.Compute(); err == nil {
		t.Error("Compute for wrong content did not fail!")
	}
}
