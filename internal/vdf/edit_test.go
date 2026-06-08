package vdf

import "testing"

func TestDocumentClone(t *testing.T) {
	doc := NewDocument(
		NewNode("root",
			NewValue("item", "one"),
			NewValue("item", "two"),
		),
	)
	clone := doc.Clone()
	clone.Path("root", "item").Value = "changed"

	if got := doc.Path("root", "item").Value; got != "one" {
		t.Fatalf("original changed after clone mutation: %q", got)
	}
	if got := clone.Path("root", "item").Value; got != "changed" {
		t.Fatalf("clone value = %q", got)
	}
}

func TestDocumentAppendSetFirstAndRemove(t *testing.T) {
	doc := NewDocument(NewValue("item", "one"), NewValue("item", "two"))

	doc.Append(NewValue("tail", "ok"))
	if got := doc.First("tail").Value; got != "ok" {
		t.Fatalf("tail = %q", got)
	}

	replaced := doc.SetFirst(NewValue("item", "replacement"))
	if !replaced {
		t.Fatalf("SetFirst should replace existing node")
	}
	items := doc.All("item")
	if len(items) != 2 {
		t.Fatalf("item count = %d", len(items))
	}
	if items[0].Value != "replacement" || items[1].Value != "two" {
		t.Fatalf("SetFirst changed duplicate semantics: %#v", items)
	}

	replaced = doc.SetFirst(NewValue("missing", "appended"))
	if replaced {
		t.Fatalf("SetFirst should append missing key")
	}
	if got := doc.First("missing").Value; got != "appended" {
		t.Fatalf("missing = %q", got)
	}

	removed := doc.RemoveFirst("item")
	if removed == nil || removed.Value != "replacement" {
		t.Fatalf("RemoveFirst removed %#v", removed)
	}
	if got := doc.First("item").Value; got != "two" {
		t.Fatalf("remaining item = %q", got)
	}

	all := doc.RemoveAll("item")
	if len(all) != 1 || all[0].Value != "two" {
		t.Fatalf("RemoveAll removed %#v", all)
	}
	if doc.First("item") != nil {
		t.Fatalf("item should be gone")
	}
}

func TestNodeAppendSetFirstAndRemove(t *testing.T) {
	node := NewValue("root", "scalar")
	node.Append(NewValue("item", "one"), NewValue("item", "two"))

	if !node.IsObject() || node.Value != "" {
		t.Fatalf("Append should convert value node to object")
	}
	if got := len(node.All("item")); got != 2 {
		t.Fatalf("item count = %d", got)
	}

	if !node.SetFirst(NewValue("item", "replacement")) {
		t.Fatalf("SetFirst should replace child")
	}
	items := node.All("item")
	if items[0].Value != "replacement" || items[1].Value != "two" {
		t.Fatalf("SetFirst changed duplicate semantics: %#v", items)
	}

	if node.SetFirst(nil) {
		t.Fatalf("SetFirst(nil) should be false")
	}
	if got := node.RemoveFirst("item").Value; got != "replacement" {
		t.Fatalf("RemoveFirst = %q", got)
	}
	if got := node.RemoveAll("item"); len(got) != 1 || got[0].Value != "two" {
		t.Fatalf("RemoveAll = %#v", got)
	}
}

func TestEditNilReceivers(t *testing.T) {
	var doc *Document
	if doc.Clone() != nil {
		t.Fatalf("nil document clone should be nil")
	}
	doc.Append(NewValue("x", "y"))
	if doc.SetFirst(NewValue("x", "y")) {
		t.Fatalf("nil document SetFirst should be false")
	}
	if doc.RemoveFirst("x") != nil {
		t.Fatalf("nil document RemoveFirst should be nil")
	}
	if doc.RemoveAll("x") != nil {
		t.Fatalf("nil document RemoveAll should be nil")
	}

	var node *Node
	if node.Clone() != nil {
		t.Fatalf("nil node clone should be nil")
	}
	node.Append(NewValue("x", "y"))
	if node.SetFirst(NewValue("x", "y")) {
		t.Fatalf("nil node SetFirst should be false")
	}
	if node.RemoveFirst("x") != nil {
		t.Fatalf("nil node RemoveFirst should be nil")
	}
	if node.RemoveAll("x") != nil {
		t.Fatalf("nil node RemoveAll should be nil")
	}
}
