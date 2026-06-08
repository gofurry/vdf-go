package vdf

import "testing"

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		`"root" { "k" "v" }`,
		`key value`,
		`"unterminated`,
		`#include "base.vdf"`,
		`"root" { "path" "C:\\Steam" }`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = ParseString(input,
			WithMaxDepth(64),
			WithMaxTokenBytes(64*1024),
			WithMaxNodes(10000),
		)
	})
}
