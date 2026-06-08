package main

import (
	"fmt"
	"log"

	vdf "github.com/gofurry/vdf-go"
)

func main() {
	doc, err := vdf.ParseString(`"AppState" { "appid" "730" "name" "Counter-Strike 2" }`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.Path("AppState", "name").AsString())
}
