package vdf

import "testing"

func TestValueHelpers(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"1", true},
		{"true", true},
		{"yes", true},
		{"on", true},
		{"0", false},
		{"false", false},
		{"no", false},
		{"off", false},
	}
	for _, tt := range tests {
		got, err := NewValue("k", tt.value).AsBool()
		if err != nil {
			t.Fatalf("AsBool(%q) error = %v", tt.value, err)
		}
		if got != tt.want {
			t.Fatalf("AsBool(%q) = %v", tt.value, got)
		}
	}

	if got, err := NewValue("k", "42").AsInt(); err != nil || got != 42 {
		t.Fatalf("AsInt() = %d, %v", got, err)
	}
	if got, err := NewValue("k", "42").AsUint64(); err != nil || got != 42 {
		t.Fatalf("AsUint64() = %d, %v", got, err)
	}
	if got, err := NewValue("k", "1.5").AsFloat64(); err != nil || got != 1.5 {
		t.Fatalf("AsFloat64() = %f, %v", got, err)
	}
	if _, err := NewValue("k", "maybe").AsBool(); err == nil {
		t.Fatalf("expected bool error")
	}
}
