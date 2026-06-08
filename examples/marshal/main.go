package main

import (
	"fmt"
	"log"

	vdf "github.com/gofurry/vdf-go"
)

func main() {
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
}
