package test

import (
	"testing"
)

func TestFlushAll(t *testing.T) {
	err := client.FlushAll(ctx)
	if err != nil {
		t.Fatalf("other FlushAll error %s", err)
	}
}