# 架构与目录结构

`vdf-go` 的目录结构目标是：保持公开导入路径稳定，同时把内部实现和测试资源整理到可扩展的位置。

## 目录布局

```text
.
├── export.go              # 根包公开 facade，保持 github.com/gofurry/vdf-go 导入路径
├── doc.go                 # package-level GoDoc
├── examples_test.go       # 根包公开 facade 示例
├── internal/
│   └── vdf/               # parser、encoder、AST、options、errors、editing 实现和测试
├── testdata/
│   ├── valid/             # 支持的 VDF / ACF / KeyValues 风格 CFG fixtures
│   ├── malformed/         # parser hardening 的坏样本
│   └── unsupported/       # 明确不支持的格式边界样本
├── examples/              # 可运行示例
└── docs/                  # 英文和中文文档
```

## 公开 API 边界

用户继续使用：

```go
import vdf "github.com/gofurry/vdf-go"
```

根包通过 `export.go` 暴露稳定 API，内部实现放在 `internal/vdf`。这样后续可以继续调整 parser、lexer、encoder 或 AST 实现，而不会改变用户导入路径。

## 为什么不把公开包移动到 `vdf/`

如果把公开包移动到 `vdf/` 子目录，用户导入路径会变成：

```go
github.com/gofurry/vdf-go/vdf
```

这会破坏当前更短、更自然的 module-root 导入路径。因此本仓库保留根包作为公开入口，只把实现细节移入 `internal/`。

## Fixture 分类策略

- `testdata/valid`：当前 parser 支持的文本格式。
- `testdata/malformed`：必须安全失败的坏输入。
- `testdata/unsupported`：用于说明边界的非目标格式。

fixture 必须脱敏，不应包含真实账号 ID、认证 token、私有安装路径或机器特定信息。

## 测试位置

实现相关测试放在 `internal/vdf`，这样 coverage 能直接反映 parser、encoder、AST 和 editing 代码。根包只保留 facade 级别的 examples，用来验证公开导入路径仍然可用。
