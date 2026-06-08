# 兼容性说明

## 当前状态

当前 `vdf-go` 的核心能力集中在文本 VDF / Valve KeyValues：

- `.vdf`
- `.acf`
- 常见 KeyValues 风格文本结构

当前 parser 已支持 quoted / unquoted key-value、嵌套对象、重复 key、顺序保留、`//` 注释、基础转义和带位置的错误。

## v0.2.x 重点：补齐 Steam `.cfg` 兼容调查

Steam / Valve 生态里确实存在大量 `.cfg` 文件，但 `.cfg` 不一定等同于 VDF。

`v0.2.x` 需要先建立兼容性矩阵，把 `.cfg` 分成至少两类：

- **KeyValues 风格 `.cfg`**：结构接近 VDF / KeyValues，可以考虑由当前 parser 支持或通过小扩展支持。
- **Source / console command 风格 `.cfg`**：更接近逐行命令，例如 bind、alias、cvar 设置，不应该直接伪装成 VDF AST。

在没有完成样本归类前，不应宣称 `vdf-go` 已完整支持所有 `.cfg`。

## 保持不变的边界

- 不自动执行 `#include` / `#base`。
- 不自动读取 directive 指向的文件。
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
