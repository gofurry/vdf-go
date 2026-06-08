# 使用说明

## 安装

```sh
go get github.com/gofurry/vdf-go
```

## 解析文本 VDF / KeyValues

```go
doc, err := vdf.ParseString(`"AppState" { "appid" "730" "name" "Counter-Strike 2" }`)
if err != nil {
    panic(err)
}

appid := doc.Path("AppState", "appid").AsString()
fmt.Println(appid)
```

## 解析文件

```go
doc, err := vdf.ParseFile("appmanifest_730.acf")
if err != nil {
    panic(err)
}

name := doc.Path("AppState", "name").AsString()
fmt.Println(name)
```

`ParseFile` 只读取传入的文件路径，不会自动展开 `#include` 或 `#base`。

## 重复 key 与顺序

VDF / KeyValues 中同级 key 可能重复，顺序也可能有意义，所以 `vdf-go` 不使用 `map` 表示对象。

```go
doc, err := vdf.ParseString(`"root" { "item" "one" "item" "two" }`)
if err != nil {
    panic(err)
}

items := doc.First("root").All("item")
for _, item := range items {
    fmt.Println(item.AsString())
}
```

`First` 返回第一个匹配节点，`All` 返回全部匹配节点，`Path` 返回路径上的第一个匹配节点。

## 构造与写出

```go
doc := vdf.NewDocument(
    vdf.NewNode("AppState",
        vdf.NewValue("appid", "730"),
        vdf.NewValue("name", "Counter-Strike 2"),
    ),
)

text, err := vdf.MarshalString(doc)
if err != nil {
    panic(err)
}

fmt.Print(text)
```

Marshal 输出是稳定、可读的 VDF，但不会保留原始注释、空行、缩进风格或 quote 风格。

## 小型 AST 操作

如果需要直接修改 `Document` / `Node`，可以使用小型操作 helper：

```go
doc.Append(vdf.NewValue("new_key", "value"))
doc.SetFirst(vdf.NewValue("name", "Updated Name"))
removed := doc.RemoveFirst("old_key")
all := doc.RemoveAll("duplicate_key")
copy := doc.Clone()
```

`SetFirst` 只替换第一个同名 key，不会删除后续重复 key。如果 key 不存在，它会追加新节点。

`RemoveFirst` 只移除第一个匹配项；`RemoveAll` 会移除所有匹配项。这个命名是有意为之，避免在保留重复 key 的 AST 中发生隐式覆盖。

`Clone` 会深拷贝文档树，适合在修改前保留原始结构。

## `.cfg` 说明

`vdf-go` 支持 KeyValues 风格 `.cfg`，但不支持 Source / console command 风格 `.cfg`。

支持范围：

- `.vdf` 文本 KeyValues
- `.acf` appmanifest
- KeyValues 风格 `.cfg`

不支持范围：

- `bind "w" "+forward"` 这类命令式 `.cfg`
- binary VDF
- `shortcuts.vdf`
- Steam 本地库扫描

详细边界见 [兼容性说明](compatibility.md)。

## Directive 与 Condition Token

`#include` 和 `#base` 默认会被忽略，并且永远不会自动读取文件。

```go
doc, err := vdf.ParseString(`#include "base.vdf" "root" { "k" "v" }`)
```

如果需要保留 directive 节点：

```go
doc, err := vdf.ParseString(input, vdf.WithPreserveDirectives(true))
```

`[$WIN32]` 这类 condition token 可以出现在 value 或 object 后面。当前策略是接受并跳过，不做平台判断，关联节点无条件保留。

## 错误处理

解析错误会返回 `*vdf.ParseError`，包含 line、column 和 offset。

```go
var parseErr *vdf.ParseError
if errors.As(err, &parseErr) {
    fmt.Println(parseErr.Line, parseErr.Column, parseErr.Offset)
}
```

## 并发说明

`Document` 和 `Node` 是普通可变结构。只读访问在没有其他 goroutine 修改同一对象时是安全的；如果需要并发修改，调用方需要自己加锁或复制数据。
