package stats

import "testing"

func TestOnlineMapRefCount(t *testing.T) {
	om := NewOnlineMap()
	om.AddIP("1.1.1.1")
	om.AddIP("1.1.1.1")
	om.AddIP("2.2.2.2")

	if got := om.Count(); got != 2 {
		t.Fatalf("count mismatch: got=%d want=2", got)
	}

	om.RemoveIP("1.1.1.1")
	if got := om.Count(); got != 2 {
		t.Fatalf("count mismatch after first remove: got=%d want=2", got)
	}

	om.RemoveIP("1.1.1.1")
	if got := om.Count(); got != 1 {
		t.Fatalf("count mismatch after second remove: got=%d want=1", got)
	}
}

func TestOnlineMapIPTimeMapReturnsSnapshot(t *testing.T) {
	om := NewOnlineMap()
	om.AddIP("3.3.3.3")

	snapshot := om.IPTimeMap()
	delete(snapshot, "3.3.3.3")

	if got := om.Count(); got != 1 {
		t.Fatalf("online map mutated through snapshot: got=%d want=1", got)
	}
}

func TestOnlineMapIPUnixMapReturnsSnapshot(t *testing.T) {
	om := NewOnlineMap()
	om.AddIP("4.4.4.4")

	snapshot := om.IPUnixMap()
	delete(snapshot, "4.4.4.4")

	if got := om.Count(); got != 1 {
		t.Fatalf("online map mutated through unix snapshot: got=%d want=1", got)
	}
}
