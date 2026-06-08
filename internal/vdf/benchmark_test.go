package vdf

import (
	"os"
	"testing"
)

var benchDoc = NewDocument(
	NewNode("AppState",
		NewValue("appid", "730"),
		NewValue("name", "Counter-Strike 2"),
		NewValue("StateFlags", "4"),
		NewNode("InstalledDepots",
			NewNode("731",
				NewValue("manifest", "123456789"),
				NewValue("size", "123456"),
			),
		),
	),
)

func BenchmarkParse(b *testing.B) {
	data, err := Marshal(benchDoc)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Parse(data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(benchDoc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseFixtures(b *testing.B) {
	fixtures := []string{
		"../../testdata/valid/libraryfolders.vdf",
		"../../testdata/valid/config.vdf",
		"../../testdata/valid/loginusers.vdf",
		"../../testdata/valid/appmanifest_730.acf",
		"../../testdata/valid/appmanifest_570.acf",
		"../../testdata/valid/sample_keyvalues.cfg",
	}
	for _, path := range fixtures {
		data, err := os.ReadFile(path)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(path, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := Parse(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkMarshalFixtures(b *testing.B) {
	fixtures := []string{
		"../../testdata/valid/libraryfolders.vdf",
		"../../testdata/valid/config.vdf",
		"../../testdata/valid/loginusers.vdf",
		"../../testdata/valid/appmanifest_730.acf",
		"../../testdata/valid/appmanifest_570.acf",
		"../../testdata/valid/sample_keyvalues.cfg",
	}
	for _, path := range fixtures {
		doc, err := ParseFile(path)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(path, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := Marshal(doc); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
