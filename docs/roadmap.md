# vdf-go v0.1.x Roadmap

> 目标：把 `vdf-go` 做成一个小而稳定的 Go 文本 VDF / Valve KeyValues 解析库。  
> 范围：本文只规划 `v0.1.x`，不承诺 `v0.2.0`、`v1.0.0` 或 binary VDF。  
> 定位：通用 VDF / KeyValues parser，不绑定 Steam Web API，不做 Steam 客户端扫描工具。

## Current Position

`vdf-go` 是一个用于解析和处理 Valve Data Format（VDF / KeyValues）文本文件的小型 Go 工具库。

核心任务是解析 Valve / Steam 生态里常见的文本 VDF / KeyValues 文件，例如：

- `libraryfolders.vdf`
- `appmanifest_*.acf`
- `config.vdf`
- `loginusers.vdf`
- Source / Valve KeyValues 风格配置文件

`v0.1.x` 阶段不做 Steam 专用库，不承担 Steam 本地库扫描、账号信息读取、binary VDF 或 shortcuts 解析。

## Roadmap Strategy

优先把文本 parser、AST、错误位置、marshal、输入安全限制和测试做稳。核心包只理解 VDF / KeyValues 语法，不理解 Steam、AppID、libraryfolders 或 appmanifest 的业务含义。

最重要的三个契约：

1. 保留顺序和重复 key。
2. 默认不做危险的 include/base 文件读取。
3. 核心模型不绑定 Steam 业务。

## Version Plan

### v0.1.0 - Parser MVP

**Status:** Completed  
**Scope:** User-facing / Developer-facing / Testing  
**Goal:** 完成最小可用文本 VDF parser。

#### Focus

- `Document` / `Node`
- parser entrypoints
- duplicated key and order preservation
- positioned parse errors

#### Tasks

- [x] Add `Document` / `Node`
- [x] Add `Parse`, `ParseString`, `ParseReader`, and `ParseFile`
- [x] Support quoted and unquoted tokens
- [x] Support nested objects and multiple root keys
- [x] Preserve duplicate keys and child order
- [x] Add `First`, `All`, and `Path`
- [x] Add `ParseError` with line / column / offset
- [x] Add basic parser tests

#### Acceptance Criteria

- Parser handles common text VDF / KeyValues syntax.
- Repeated sibling keys are preserved.
- Parser errors include source position.

---

### v0.1.1 - Marshal and Constructors

**Status:** Completed  
**Scope:** User-facing / Testing  
**Goal:** 补充可写能力和构造 API。

#### Focus

- stable readable VDF output
- document and node construction
- roundtrip behavior

#### Tasks

- [x] Add `Marshal`, `MarshalString`, and `Write`
- [x] Add `NewDocument`, `NewNode`, and `NewValue`
- [x] Add encode options
- [x] Add marshal roundtrip tests

#### Acceptance Criteria

- Documents can be constructed in Go and encoded as readable VDF.
- Marshal does not sort keys by default.
- Marshal tests cover escaping, nested objects, duplicate keys, and roundtrip parsing.

---

### v0.1.2 - Value Helpers and Safety Limits

**Status:** Completed  
**Scope:** User-facing / Security/Safety / Testing  
**Goal:** 提升使用体验和输入安全。

#### Focus

- string-to-type helpers
- parser limits
- conservative defaults

#### Tasks

- [x] Add `AsString`, `AsInt`, `AsUint64`, `AsFloat64`, and `AsBool`
- [x] Add `MaxDepth`, `MaxTokenBytes`, and `MaxNodes`
- [x] Add limit error tests
- [x] Document bool helper semantics

#### Acceptance Criteria

- Type helpers return errors instead of panics.
- Extreme depth, token size, and node count are bounded by default.
- Limit failures return positioned parse errors.

---

### v0.1.3 - Real Steam File Fixtures

**Status:** Completed  
**Scope:** Compatibility / Testing  
**Goal:** 验证真实 Steam 文本文件兼容性。

#### Focus

- `libraryfolders.vdf`
- `appmanifest_*.acf`
- malformed fixtures

#### Tasks

- [x] Add `libraryfolders.vdf` fixture
- [x] Add `appmanifest_730.acf` fixture
- [x] Add Windows path tests
- [x] Add duplicate key tests
- [x] Add escape sequence tests
- [x] Add malformed fixture tests

#### Acceptance Criteria

- Common Steam text VDF shapes parse successfully.
- Fixtures validate structure only, not Steam business semantics.
- Malformed fixtures fail cleanly.

---

### v0.1.4 - Documentation and Examples

**Status:** Completed  
**Scope:** Documentation / User-facing  
**Goal:** 完善开源使用体验。

#### Focus

- README
- examples
- compatibility and security docs

#### Tasks

- [x] Add README quick start
- [x] Add `examples/parse`
- [x] Add `examples/marshal`
- [x] Add `docs/compatibility.md`
- [x] Add `docs/roadmap.md`
- [x] Add `SECURITY.md`
- [x] Add `CONTRIBUTING.md`

#### Acceptance Criteria

- New users can parse and marshal from README snippets.
- Compatibility boundary is documented.
- Unsafe include/base expansion is explicitly out of scope.

---

### v0.1.5 - Fuzz and Hardening

**Status:** Completed  
**Scope:** Stability / Testing / CI  
**Goal:** 增强 parser 稳定性和 release 基线。

#### Focus

- fuzzing
- benchmarks
- CI

#### Tasks

- [x] Add `FuzzParse`
- [x] Add malformed input seed corpus
- [x] Add parser benchmark
- [x] Add marshal benchmark
- [x] Add CI for `gofmt`, `go vet`, `go test`, and `staticcheck`

#### Acceptance Criteria

- Parser fuzzing treats malformed data as ordinary errors.
- Benchmarks cover parser and marshal paths.
- CI validates formatting, vet, tests, and static analysis.

## v0.1.x Non-Goals

The full `v0.1.x` line does not implement:

- binary VDF
- `shortcuts.vdf`
- Steam library scanner
- appmanifest strongly typed business model
- `libraryfolders` strongly typed business model
- struct decoder
- comments round-trip
- include/base automatic file reading
- formatter CLI
- `steam-go` bridge layer
- v1.0 stability promise

## Acceptance Checklist

- [x] 能解析常见文本 VDF / KeyValues 文件
- [x] 能解析 `.acf` appmanifest
- [x] 能解析 `libraryfolders.vdf`
- [x] 保留重复 key
- [x] 保留节点顺序
- [x] 错误包含 line / column / offset
- [x] 有 marshal 能力
- [x] 有基础类型转换 helper
- [x] 有输入限制，避免极端文件导致资源耗尽
- [x] 有 fuzz 测试
- [x] 有 README、examples、roadmap、compatibility docs
- [x] CI 至少包含 `go test ./...`、`go vet ./...`、`staticcheck`

## Relationship with steam-go

`v0.1.x` 阶段，`vdf-go` 不需要和 `steam-go` 强绑定。

推荐关系：

```text
vdf-go
  通用 VDF / KeyValues parser

steam-go/addons/vdf
  后续桥接层
  只 re-export 常用类型或提供 Steam 本地文件 helper

SteamScope / GoFurry tools
  未来消费 vdf-go 或 steam-go/addons/vdf
```

等 `vdf-go` parser 稳定后，再考虑后续 Steam helper、`steam-go/addons/vdf`、binary VDF 或 SteamScope 本地库扫描。
