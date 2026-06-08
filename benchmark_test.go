package vdf

import "testing"

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
