package vdf_test

import (
	"fmt"
	"log"

	vdf "github.com/gofurry/vdf-go"
)

func ExampleParseString() {
	doc, err := vdf.ParseString(`"AppState" { "appid" "730" "name" "Counter-Strike 2" }`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.Path("AppState", "appid").AsString())
	// Output:
	// 730
}

func ExampleNode_All() {
	doc, err := vdf.ParseString(`"root" { "item" "one" "item" "two" }`)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range doc.First("root").All("item") {
		fmt.Println(item.AsString())
	}
	// Output:
	// one
	// two
}

func ExampleMarshalString() {
	doc := vdf.NewDocument(
		vdf.NewNode("AppState",
			vdf.NewValue("appid", "730"),
			vdf.NewValue("name", "Counter-Strike 2"),
		),
	)

	text, err := vdf.MarshalString(doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(text)
	// Output:
	// "AppState"
	// {
	// 	"appid"	"730"
	// 	"name"	"Counter-Strike 2"
	// }
}
