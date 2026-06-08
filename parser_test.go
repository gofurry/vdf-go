package vdf

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestParseBasicFeatures(t *testing.T) {
	input := `"AppState"
{
	"appid"	"730"
	name "Counter-Strike 2"
	"installdir" "Counter-Strike Global Offensive"
	"windows_path" "C:\\Program Files (x86)\\Steam"
	// comment
	"escaped" "line\nquote\"tab\tbackslash\\"
}
"root2" "value2"
`
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	app := doc.First("AppState")
	if app == nil || !app.IsObject() {
		t.Fatalf("AppState object not found")
	}
	if got := app.First("appid").Value; got != "730" {
		t.Fatalf("appid = %q", got)
	}
	if got := app.First("windows_path").Value; got != `C:\Program Files (x86)\Steam` {
		t.Fatalf("windows_path = %q", got)
	}
	if got := app.First("escaped").Value; got != "line\nquote\"tab\tbackslash\\" {
		t.Fatalf("escaped = %q", got)
	}
	if got := doc.First("root2").Value; got != "value2" {
		t.Fatalf("root2 = %q", got)
	}
}

func TestParseDuplicateKeysAndPath(t *testing.T) {
	doc, err := ParseString(`"root" { "item" "one" "item" "two" "child" { "leaf" "ok" } }`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	root := doc.First("root")
	items := root.All("item")
	if len(items) != 2 {
		t.Fatalf("duplicate item count = %d", len(items))
	}
	if items[0].Value != "one" || items[1].Value != "two" {
		t.Fatalf("items out of order: %#v", items)
	}
	if got := doc.Path("root", "child", "leaf").Value; got != "ok" {
		t.Fatalf("path value = %q", got)
	}
}

func TestParseEmptyValueAndEmptyObject(t *testing.T) {
	doc, err := ParseString(`"root" { "empty_value" "" "empty_object" { } }`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	root := doc.First("root")
	if !root.First("empty_value").IsValue() {
		t.Fatalf("empty_value should be value")
	}
	if !root.First("empty_object").IsObject() {
		t.Fatalf("empty_object should be object")
	}
}

func TestDirectivesDefaultIgnoredAndPreserved(t *testing.T) {
	input := `#include "base.vdf" "root" { #base "defaults.vdf" "k" "v" }`
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if doc.First("#include") != nil {
		t.Fatalf("directive should be ignored by default")
	}
	if doc.First("root").First("#base") != nil {
		t.Fatalf("nested directive should be ignored by default")
	}

	doc, err = ParseString(input, WithPreserveDirectives(true))
	if err != nil {
		t.Fatalf("ParseString(preserve) error = %v", err)
	}
	if got := doc.First("#include").Value; got != "base.vdf" {
		t.Fatalf("include = %q", got)
	}
	if got := doc.First("root").First("#base").Value; got != "defaults.vdf" {
		t.Fatalf("base = %q", got)
	}
}

func TestParseOptions(t *testing.T) {
	if _, err := ParseString(`key value`, WithBareTokens(false)); err == nil {
		t.Fatalf("expected bare token error")
	}
	doc, err := ParseString(`"k" "v" // comment`, WithComments(false))
	if err != nil {
		t.Fatalf("ParseString(comments disabled) error = %v", err)
	}
	if got := doc.First("//").Value; got != "comment" {
		t.Fatalf("disabled comment token = %q", got)
	}
	doc, err = ParseString(`"k" "a\nb"`, WithEscapeSequences(false))
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if got := doc.First("k").Value; got != `a\nb` {
		t.Fatalf("escape-disabled value = %q", got)
	}
}

func TestParseErrorsHavePosition(t *testing.T) {
	tests := []string{
		`"root" { "k" "v"`,
		`"root" { "k" "v" } }`,
		`"root" { "k" "unterminated`,
		`"root" { "k" }`,
	}
	for _, input := range tests {
		_, err := ParseString(input)
		if err == nil {
			t.Fatalf("ParseString(%q) succeeded", input)
		}
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("error %T is not *ParseError", err)
		}
		if parseErr.Line <= 0 || parseErr.Column <= 0 || parseErr.Offset < 0 {
			t.Fatalf("bad position: %#v", parseErr)
		}
	}
}

func TestParseLimits(t *testing.T) {
	if _, err := ParseString(`"a" { "b" { "c" "d" } }`, WithMaxDepth(1)); err == nil {
		t.Fatalf("expected max depth error")
	}
	if _, err := ParseString(`"long" "value"`, WithMaxTokenBytes(3)); err == nil {
		t.Fatalf("expected max token error")
	}
	if _, err := ParseString(`"a" "1" "b" "2"`, WithMaxNodes(1)); err == nil {
		t.Fatalf("expected max nodes error")
	}
}

func TestParseReader(t *testing.T) {
	doc, err := ParseReader(strings.NewReader(`"k" "v"`))
	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}
	if got := doc.First("k").Value; got != "v" {
		t.Fatalf("value = %q", got)
	}
}

func TestFixtures(t *testing.T) {
	tests := []struct {
		path string
		key  string
	}{
		{"testdata/simple.vdf", "root"},
		{"testdata/nested.vdf", "root"},
		{"testdata/duplicate_keys.vdf", "root"},
		{"testdata/libraryfolders.vdf", "libraryfolders"},
		{"testdata/appmanifest_730.acf", "AppState"},
	}
	for _, tt := range tests {
		doc, err := ParseFile(tt.path)
		if err != nil {
			t.Fatalf("ParseFile(%q) error = %v", tt.path, err)
		}
		if doc.First(tt.key) == nil {
			t.Fatalf("ParseFile(%q) missing %q", tt.path, tt.key)
		}
	}

	doc, err := ParseFile("testdata/duplicate_keys.vdf")
	if err != nil {
		t.Fatalf("ParseFile(duplicate) error = %v", err)
	}
	if got := len(doc.First("root").All("item")); got != 3 {
		t.Fatalf("fixture duplicate count = %d", got)
	}

	app, err := ParseFile("testdata/appmanifest_730.acf")
	if err != nil {
		t.Fatalf("ParseFile(appmanifest) error = %v", err)
	}
	if got := app.Path("AppState", "InstalledDepots", "731", "manifest").Value; got == "" {
		t.Fatalf("missing depot manifest")
	}
}

func TestMalformedFixtures(t *testing.T) {
	for _, path := range []string{
		"testdata/malformed_missing_brace.vdf",
		"testdata/malformed_unterminated_quote.vdf",
	} {
		if _, err := ParseFile(path); err == nil {
			t.Fatalf("ParseFile(%q) succeeded", path)
		}
	}
}

func TestParseFileReadError(t *testing.T) {
	path := "testdata/does-not-exist.vdf"
	if _, err := ParseFile(path); err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ParseFile(%q) error = %v", path, err)
	}
}
