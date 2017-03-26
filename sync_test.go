package main

import "testing"

func TestSync(t *testing.T) {
	err := SyncTemplates()
	if err != nil {
		t.Errorf("sync err: %v\n", err)
	}
}
