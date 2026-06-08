package vdf

import (
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	doc := NewDocument(
		NewNode("AppState",
			NewValue("appid", "730"),
			NewValue("name", "Counter-Strike 2"),
			NewNode("empty"),
		),
		NewValue("root2", "line\nquote\"slash\\"),
	)

	got, err := MarshalString(doc)
	if err != nil {
		t.Fatalf("MarshalString() error = %v", err)
	}
	want := "\"AppState\"\n{\n\t\"appid\"\t\"730\"\n\t\"name\"\t\"Counter-Strike 2\"\n\t\"empty\"\n\t{\n\t}\n}\n\"root2\"\t\"line\\nquote\\\"slash\\\\\"\n"
	if got != want {
		t.Fatalf("MarshalString() =\n%q\nwant\n%q", got, want)
	}

	roundtrip, err := ParseString(got)
	if err != nil {
		t.Fatalf("roundtrip parse error = %v", err)
	}
	if got := roundtrip.Path("AppState", "name").Value; got != "Counter-Strike 2" {
		t.Fatalf("roundtrip name = %q", got)
	}
}

func TestMarshalSortKeysStable(t *testing.T) {
	doc := NewDocument(
		NewValue("b", "1"),
		NewValue("a", "first"),
		NewValue("a", "second"),
	)
	got, err := MarshalString(doc, WithSortKeys(true), WithQuoteKeys(false), WithQuoteValues(false))
	if err != nil {
		t.Fatalf("MarshalString() error = %v", err)
	}
	want := "a\tfirst\na\tsecond\nb\t1\n"
	if got != want {
		t.Fatalf("MarshalString() = %q, want %q", got, want)
	}
}

func TestWriteErrors(t *testing.T) {
	if err := Write(nil, NewDocument()); err == nil {
		t.Fatalf("expected nil writer error")
	}
	var b strings.Builder
	if err := Write(&b, nil); err == nil {
		t.Fatalf("expected nil document error")
	}
}
