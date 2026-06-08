# vdf-go

A small Go parser and toolkit for Valve Data Format (VDF / KeyValues) text files.

`vdf-go` focuses on the generic text format used by Valve and Steam files such as
`libraryfolders.vdf`, `appmanifest_*.acf`, `config.vdf`, and Source-style
KeyValues files, including KeyValues-style `.cfg` files. It is not a Steam
client scanner and does not implement binary VDF.

## Install

```sh
go get github.com/gofurry/vdf-go
```

## Quick Start

```go
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

	fmt.Println(doc.Path("AppState", "appid").AsString())
}
```

## Core Model

VDF can contain multiple root keys, repeated sibling keys, and meaningful order.
For that reason, `vdf-go` uses slices instead of maps:

```go
type Document struct {
	Nodes []*Node
}

type Node struct {
	Key      string
	Value    string
	Children []*Node
}
```

`Children != nil` means the node is an object. `Children == nil` means the node
is a scalar value. Use `IsObject()` and `IsValue()` when the distinction matters.

## Parsing

```go
doc, err := vdf.Parse(data)
doc, err := vdf.ParseString(text)
doc, err := vdf.ParseReader(reader)
doc, err := vdf.ParseFile("appmanifest_730.acf")
```

The parser supports quoted and unquoted keys/values, nested objects, multiple
root keys, repeated keys, `//` comments, and common escape sequences.

Parser options include:

```go
doc, err := vdf.ParseString(input,
	vdf.WithMaxDepth(128),
	vdf.WithMaxTokenBytes(1<<20),
	vdf.WithMaxNodes(1_000_000),
)
```

`#include` and `#base` are never executed and never read files. By default they
are ignored. Use `WithPreserveDirectives(true)` to keep them as ordinary
key/value nodes. Other `#...` keys are parsed as ordinary keys.

Condition tokens such as `[$WIN32]` are accepted after a value or object and
discarded. `vdf-go` does not evaluate platform conditions; the associated node
is parsed unconditionally.

## Querying

```go
root := doc.First("AppState")
allItems := root.All("item")
appid := doc.Path("AppState", "appid")
```

`First` returns the first matching node. `All` returns every matching node.
Duplicate keys are preserved and never silently overwritten.

## Writing

```go
doc := vdf.NewDocument(
	vdf.NewNode("AppState",
		vdf.NewValue("appid", "730"),
		vdf.NewValue("name", "Counter-Strike 2"),
	),
)

text, err := vdf.MarshalString(doc)
```

Marshal output is stable and readable, but it does not preserve original
comments, blank lines, or formatting style.

## Value Helpers

VDF values are strings. Helpers are provided for convenience:

```go
node.AsString()
node.AsInt()
node.AsUint64()
node.AsFloat64()
node.AsBool()
```

Boolean parsing accepts `1/true/yes/on` and `0/false/no/off`. This is a
`vdf-go` convenience behavior, not an official Valve type system.

## Editing

Small tree editing helpers are available for callers that want to modify the AST
without converting it to a map:

```go
doc.Append(vdf.NewValue("new_key", "value"))
doc.SetFirst(vdf.NewValue("name", "Updated Name"))
removed := doc.RemoveFirst("old_key")
all := doc.RemoveAll("duplicate_key")
copy := doc.Clone()
```

`SetFirst` only replaces the first matching key. It does not remove later
duplicates.

## Errors

Parse failures return `*vdf.ParseError` with line, column, and byte offset:

```go
var parseErr *vdf.ParseError
if errors.As(err, &parseErr) {
	fmt.Println(parseErr.Line, parseErr.Column, parseErr.Offset)
}
```

## Concurrency Safety

`Document` and `Node` are plain mutable structs. Read-only access is safe when no
goroutine mutates the same values. If callers mutate documents or nodes, they are
responsible for synchronization.

Parsing and marshaling functions do not use package-level mutable state.

## Compatibility

`vdf-go` supports text VDF / KeyValues only, including KeyValues-style `.cfg`
files. It does not support Source / console command-style `.cfg`, binary VDF,
`shortcuts.vdf`, struct decoding, Steam library scanning, comment
round-tripping, or automatic include/base expansion.

See [docs/compatibility.md](docs/compatibility.md) for details.

## Chinese Documentation

Chinese documentation is maintained under [docs/zh](docs/zh/README.md),
including the active roadmap and Steam / Valve format compatibility notes.

Release preparation steps are documented in
[docs/release-checklist.md](docs/release-checklist.md).
