# 兼容性说明

## 当前状态

当前 `vdf-go` 的核心能力集中在文本 VDF / Valve KeyValues：

- `.vdf`
- `.acf`
- KeyValues 风格 `.cfg`
- 常见 KeyValues 风格文本结构

当前 parser 已支持 quoted / unquoted key-value、嵌套对象、重复 key、顺序保留、`//` 注释、基础转义和带位置的错误。

## 格式支持矩阵

| 格式 | 当前状态 | 说明 |
| --- | --- | --- |
| `.vdf` 文本 KeyValues | 支持 | 例如 `libraryfolders.vdf`、`config.vdf`。 |
| `.acf` appmanifest | 支持 | 作为文本 KeyValues 解析，不提供 Steam 业务模型。 |
| KeyValues 风格 `.cfg` | 支持 | 结构接近 VDF / KeyValues，可以由当前 parser 解析。 |
| condition token，例如 `[$WIN32]` | 接受但不求值 | 可出现在 value 或 object 后面；parser 会跳过 token，节点无条件保留。 |
| Source / console command 风格 `.cfg` | 不支持 | 逐行命令格式，例如 `bind`、`alias`、`con_enable`，不应伪装成 VDF AST。 |
| binary VDF / `shortcuts.vdf` | 不支持 | 属于不同格式，不能混入当前文本 parser。 |

## `.cfg` 分类结论

Steam / Valve 生态里存在大量 `.cfg` 文件，但 `.cfg` 不是单一格式。

- **KeyValues 风格 `.cfg`**：属于当前 `vdf-go` 的支持范围。测试样本见 `testdata/sample_keyvalues.cfg`。
- **Source / console command 风格 `.cfg`**：不属于当前 VDF / KeyValues parser 范围。样本见 `testdata/source_commands.cfg`，仅用于说明边界，不作为 parse 成功 fixture。

因此，文档中应避免写“完整支持 `.cfg`”。更准确的说法是：

> 支持 KeyValues 风格 `.cfg`；不支持 Source / console command 风格 `.cfg`。

## 已覆盖样本

当前测试样本覆盖：

- `testdata/libraryfolders.vdf`
- `testdata/config.vdf`
- `testdata/loginusers.vdf`
- `testdata/appmanifest_730.acf`
- `testdata/appmanifest_570.acf`
- `testdata/sample_keyvalues.cfg`
- `testdata/source_commands.cfg`，仅作为不支持边界样本

## 保持不变的边界

- 只有 `#include` / `#base` 被当作特殊 directive；其他 `#xxx` 作为普通 key 解析。
- 不自动执行 `#include` / `#base`。
- 不自动读取 directive 指向的文件。
- 不对 `[$WIN32]` 等 condition token 做平台判断；condition token 被接受并跳过，关联节点无条件保留。
- 不把核心包变成 Steam 本地扫描器。
- 不默认引入 Steam 业务模型。
- 不把 binary VDF / `shortcuts.vdf` 混入文本 parser 的核心路径。

## 兼容性验证方式

后续新增格式支持时，需要同时补充：

- 真实或脱敏 fixture；
- malformed fixture；
- parse 成功 / 失败用例；
- marshal 是否适用的说明；
- README 或中文文档里的边界说明。
