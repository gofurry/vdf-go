# vdf-go

[English README](../../README.md)

![License](https://img.shields.io/badge/License-MIT-6C757D?style=flat&color=3B82F6)
![Release](https://img.shields.io/github/v/release/gofurry/vdf-go?style=flat&color=blue)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go&logoColor=white)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofurry/vdf-go)](https://goreportcard.com/report/github.com/gofurry/vdf-go)

一个用于解析和写出 Valve Data Format（VDF / KeyValues）文本文件的小型 Go 库。

`vdf-go` 关注 Valve 和 Steam 生态中的通用文本格式，例如 `libraryfolders.vdf`、`appmanifest_*.acf`、`config.vdf`、Source 风格 KeyValues 文件，以及 KeyValues 风格 `.cfg` 文件。它不是 Steam 客户端扫描器，也不实现 binary VDF。

## 安装

```sh
go get github.com/gofurry/vdf-go
```

## 快速开始

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

## 核心模型

VDF 可以包含多个 root key、重复的同级 key，并且顺序可能有意义。因此 `vdf-go` 使用 slice，而不是 map：

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

`Children != nil` 表示 object node。`Children == nil` 表示 scalar value node。需要区分时使用 `IsObject()` 和 `IsValue()`。

## 解析

```go
doc, err := vdf.Parse(data)
doc, err := vdf.ParseString(text)
doc, err := vdf.ParseReader(reader)
doc, err := vdf.ParseReaderLimit(reader, 1<<20)
doc, err := vdf.ParseFile("appmanifest_730.acf")
```

parser 支持 quoted / unquoted key-value、嵌套 object、多个 root key、重复 key、`//` 注释和常见转义。

常用 parser options：

```go
doc, err := vdf.ParseString(input,
	vdf.WithMaxDepth(128),
	vdf.WithMaxTokenBytes(1<<20),
	vdf.WithMaxNodes(1_000_000),
)
```

`ParseReader` 会先把 reader 完整读入内存再解析。如果输入来自不可信来源或可能很大，请使用 `ParseReaderLimit`，或在调用前用 `io.LimitReader` 包装 reader。

`#include` 和 `#base` 永远不会被执行，也不会自动读取文件。默认会忽略它们；如果需要保留为普通 key/value node，可以使用 `WithPreserveDirectives(true)`。其他 `#...` key 会作为普通 key 解析。

`[$WIN32]` 这类 condition token 可以出现在 value 或 object 后面，并会被丢弃。`vdf-go` 不做平台条件求值，关联节点会无条件保留。

## 查询

```go
root := doc.First("AppState")
allItems := root.All("item")
appid := doc.Path("AppState", "appid")
```

`First` 返回第一个匹配节点。`All` 返回全部匹配节点。重复 key 会被保留，不会被静默覆盖。

## 写出

```go
doc := vdf.NewDocument(
	vdf.NewNode("AppState",
		vdf.NewValue("appid", "730"),
		vdf.NewValue("name", "Counter-Strike 2"),
	),
)

text, err := vdf.MarshalString(doc)
```

Marshal 输出稳定、可读，但不会保留原始注释、空行或格式风格。

## 类型转换 Helper

VDF value 本质上是字符串。库提供一些便利转换：

```go
node.AsString()
node.AsInt()
node.AsUint64()
node.AsFloat64()
node.AsBool()
```

bool 转换支持 `1/true/yes/on` 和 `0/false/no/off`。这是 `vdf-go` 的便利语义，不代表 Valve 官方类型系统。

## 编辑

如果需要修改 AST，可以使用小型树操作 helper，而不必把结构转换成 map：

```go
doc.Append(vdf.NewValue("new_key", "value"))
doc.SetFirst(vdf.NewValue("name", "Updated Name"))
removed := doc.RemoveFirst("old_key")
all := doc.RemoveAll("duplicate_key")
copy := doc.Clone()
```

`SetFirst` 只替换第一个匹配 key，不会删除后续重复 key。

## 错误

解析失败会返回带 line、column 和 byte offset 的 `*vdf.ParseError`：

```go
var parseErr *vdf.ParseError
if errors.As(err, &parseErr) {
	fmt.Println(parseErr.Line, parseErr.Column, parseErr.Offset)
}
```

## 并发安全

`Document` 和 `Node` 是普通可变结构。在没有 goroutine 同时修改同一对象时，只读访问是安全的。如果调用方需要并发修改，需要自行同步。

parse 和 marshal 函数不使用 package-level mutable state。

## 兼容性

`vdf-go` 只支持文本 VDF / KeyValues，包括 KeyValues 风格 `.cfg` 文件。不支持 Source / console command 风格 `.cfg`、binary VDF、`shortcuts.vdf`、struct decoder、Steam library scanner、comment round-trip 或自动 include/base 展开。

更多细节见 [兼容性说明](compatibility.md)。

## 文档

- [使用说明](usage.md)
- [兼容性说明](compatibility.md)
- [架构与目录结构](architecture.md)
- [路线图](roadmap.md)
- [发布检查清单](../release-checklist.md)
- [Valve KeyValues 官方文档](https://developer.valvesoftware.com/wiki/KeyValues)

## 项目结构

公开导入路径仍然是 module root：

```go
import vdf "github.com/gofurry/vdf-go"
```

实现文件位于 `internal/vdf`，根包通过小型稳定 facade 暴露 API。Fixtures 按 `testdata/valid/steam-client`、`testdata/valid/appmanifest`、`testdata/valid/source-keyvalues`、`testdata/valid/directives`、`testdata/valid/conditions`、`testdata/malformed` 和 `testdata/unsupported` 分类。

详细说明见 [架构与目录结构](architecture.md)。
