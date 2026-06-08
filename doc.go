// Package vdf parses and writes Valve Data Format (VDF / KeyValues) text files.
//
// The package is intentionally small and format-focused. It supports text VDF
// and KeyValues-style files such as .vdf, .acf, and KeyValues-style .cfg files.
// It does not scan Steam installations, decode Steam business models, parse
// binary VDF, or execute #include/#base directives.
//
// VDF can contain duplicate sibling keys and meaningful child order, so the
// core model uses slices instead of maps:
//
//	doc, err := vdf.ParseString(`"AppState" { "appid" "730" }`)
//	if err != nil {
//		// handle error
//	}
//	appid := doc.Path("AppState", "appid").AsString()
//
// Documents and nodes are plain mutable structs. Read-only access is safe when
// callers do not concurrently mutate the same values.
package vdf
