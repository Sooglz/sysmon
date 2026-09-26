package collector

import "testing"

func TestCollect(t *testing.T) {
	m, err := Collect("/")
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}
	if m.Timestamp == 0 {
		t.Fatal("timestamp is zero")
	}
	if m.MemoryPercent < 0 || m.MemoryPercent > 100 {
		t.Fatalf("memory percent out of range: %v", m.MemoryPercent)
	}
	if m.DiskPercent < 0 || m.DiskPercent > 100 {
		t.Fatalf("disk percent out of range: %v", m.DiskPercent)
	}
}
