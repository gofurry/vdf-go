# vdf-go v0.2.x 路线图

## 当前定位

`vdf-go` 当前已经完成文本 VDF / Valve KeyValues parser 的基础闭环：可以解析和写出 `.vdf` / `.acf` 这类常见文本 KeyValues 文件，保留重复 key 和节点顺序，并提供基础查询、构造、类型转换、错误位置、安全限制、fixtures、fuzz、benchmark 和 CI。

接下来 `v0.2.x` 不急着进入 `v1.0.0`，重点是扩大真实世界兼容性，把 Steam / Valve 生态里更多文本配置样本纳入验证，尤其是 `.cfg` 文件。

## 路线策略

优先顺序：

1. 先建立 `.cfg` 与更多真实文件的兼容性矩阵。
2. 再决定哪些 `.cfg` 属于 KeyValues parser 范围，哪些应留给独立 parser 或后续子包。
3. 补齐 GoDoc、CHANGELOG、中文文档和 release hygiene。
4. 用更长 fuzz、更多 malformed corpus 和真实 fixtures 稳住 parser。
5. 谨慎评估小型操作 API，避免破坏重复 key 和顺序保留契约。

`v0.2.x` 仍然保持核心包小而通用，不绑定 Steam Web API，不做 Steam 本地扫描工具。

## Version Plan

### v0.2.0 - Steam 文本样本兼容矩阵

**Status:** Completed  
**Scope:** Compatibility / Testing / Documentation  
**Goal:** 明确 `vdf-go` 对 `.vdf`、`.acf`、`.cfg` 和其他 Steam 文本配置文件的真实兼容边界。

#### Focus

- `.cfg` 文件归类
- 真实或脱敏 fixture
- 兼容性文档

#### Tasks

- [x] 收集并脱敏 Steam / Valve 生态里的代表性 `.cfg` 样本。
- [x] 区分 KeyValues 风格 `.cfg` 与 Source / console command 风格 `.cfg`。
- [x] 为 KeyValues 风格 `.cfg` 增加 parse fixtures 和回归测试。
- [x] 为命令式 `.cfg` 增加“不属于当前 VDF parser 范围”的文档说明。
- [x] 扩充 `.vdf` / `.acf` fixture，包括 `config.vdf`、`loginusers.vdf`、更多 `appmanifest_*.acf`。
- [x] 更新 `docs/zh/compatibility.md`，维护格式支持矩阵。

#### Acceptance Criteria

- 文档明确说明 `.cfg` 不是单一格式。
- KeyValues 风格 `.cfg` 有测试样本覆盖。
- 命令式 `.cfg` 不会被错误宣传为已完整支持。
- 新增 fixtures 不包含用户隐私、账号 token、安装路径中的敏感信息。

#### Notes

`v0.2.0` 只确认 `.cfg` 兼容边界：KeyValues 风格 `.cfg` 属于当前 parser 范围；Source / console command 风格 `.cfg` 暂不支持，后续如需支持应设计独立 parser 或子包。

---

### v0.2.1 - Text KeyValues 兼容性增强

**Status:** Completed  
**Scope:** User-facing / Compatibility / Parser  
**Goal:** 在不扩大 Steam 业务边界的前提下，提高 Valve / Source KeyValues 文本兼容性。

#### Focus

- directive 语义收窄
- condition token 策略
- 边缘语法回归

#### Tasks

- [x] 明确只对 `#include` 和 `#base` 做特殊 directive 识别，其他 `#xxx` 保持可预测行为。
- [x] 继续保证 directive 默认不执行、不读文件。
- [x] 设计 condition token `[$WIN32]` 的保留策略，先接受并跳过，不做条件求值。
- [x] 增加 CRLF、UTF-8、空 key、空 object、奇怪转义、行尾注释边界测试。
- [x] 增加 parser 错误消息回归测试，避免错误位置退化。

#### Acceptance Criteria

- directive 行为在 README 和中文兼容性文档中一致。
- 条件 token 有明确策略，不会被误解为已执行条件求值。
- 新增边缘语法测试通过。

#### Notes

`v0.2.1` 的 condition token 策略是“接受但不求值”：例如 `[$WIN32]` 可以出现在 value 或 object 后面，parser 会跳过该 token，关联节点仍无条件保留。后续如需条件求值，应先设计独立选项和平台上下文，不能作为默认行为。

---

### v0.2.2 - GoDoc、中文文档与发布卫生

**Status:** Completed  
**Scope:** Documentation / CI/Release / Developer-facing  
**Goal:** 把库从“能用”提升到“适合公开试用和持续维护”。

#### Focus

- package-level GoDoc
- CHANGELOG
- 中文文档目录
- release checklist

#### Tasks

- [x] 增加 `doc.go`，提供 package-level GoDoc 和最小示例。
- [x] 增加 GoDoc example tests，覆盖 parse、query、marshal。
- [x] 增加 `CHANGELOG.md`，记录 `v0.1.x` 基础能力和 `v0.2.x` 计划。
- [x] 在 `docs/zh/` 持续维护中文 roadmap、兼容性说明和使用说明。
- [x] 增加 release checklist，明确 tag、测试、fuzz、文档更新步骤。
- [x] 更新 README，链接中文文档入口。

#### Acceptance Criteria

- pkg.go.dev 第一屏能清楚说明库的定位和边界。
- 中文文档集中在 `docs/zh/`。
- release 前检查项明确，不依赖口头记忆。

#### Notes

`v0.2.2` 补齐了 `doc.go`、GoDoc example tests、`CHANGELOG.md`、中文使用说明和 release checklist。后续 public API 变化都应同步更新 GoDoc、README、中文文档和 changelog。

---

### v0.2.3 - Fuzz、Malformed Corpus 与真实样本硬化

**Status:** Completed  
**Scope:** Stability / Testing / Security/Safety  
**Goal:** 用更长 fuzz 和更多坏样本降低 parser 在真实输入下的崩溃、死循环和资源风险。

#### Focus

- longer fuzz
- malformed corpus
- resource limits

#### Tasks

- [x] 运行更长时间 fuzz，例如 `go test -run=FuzzParse -fuzz=FuzzParse -fuzztime=1m` 或更长。
- [x] 将有价值的 fuzz 发现固化为普通回归测试。
- [x] 扩充 malformed fixtures：未闭合 quote、缺失 brace、孤立 brace、超长 token、深层嵌套、directive 边界。
- [x] 增加 benchmark 对比，观察真实 fixture 下 parse / marshal 分配变化。
- [x] 检查错误消息不会输出大段原始文件内容或敏感配置片段。

#### Acceptance Criteria

- malformed 输入只返回错误，不 panic。
- fuzz 结果中有价值样本被转化为稳定测试。
- parser 限制行为有测试保护。

#### Notes

`v0.2.3` 增加了 malformed corpus、资源限制回归测试、错误消息脱敏测试，以及真实 fixture parse / marshal benchmark。已执行 1 分钟 fuzz；本轮未发现需要额外固化的新失败样本。

---

### v0.2.4 - 小型文档树操作 API 评估

**Status:** Planned  
**Scope:** API / Developer-facing  
**Goal:** 评估是否需要增加少量 AST 操作 helper，同时保护重复 key 和顺序语义。

#### Focus

- mutation helper design
- clone behavior
- compatibility risk

#### Tasks

- [ ] 评估 `Append`、`Remove`、`SetFirst`、`Clone` 等 helper 是否确实能减少用户代码复杂度。
- [ ] 明确任何 `Set` 类 API 不得默认覆盖所有重复 key。
- [ ] 如果增加 API，优先提供小而可预测的方法，并补充迁移说明。
- [ ] 如果风险大，仅在文档中给出推荐操作模式，不急于加入公共 API。

#### Acceptance Criteria

- 不破坏 `First` / `All` / `Path` 的既有语义。
- 不引入 map 化 AST。
- API 增量有测试和文档说明。

## v0.2.x 延后事项

以下内容可以继续讨论，但不作为 `v0.2.x` 必做目标：

- binary VDF。
- `shortcuts.vdf`。
- Steam 本地库扫描。
- appmanifest / libraryfolders 强类型业务模型。
- struct decoder。
- comments round-trip。
- formatter CLI。
- `steam-go` bridge layer。
- `v1.0.0` 稳定承诺。

如果后续要做 Steam helper，建议放在独立子包或独立仓库中，避免污染核心 VDF / KeyValues parser。

## v0.2.x 总体验收标准

- [x] `.cfg` 兼容边界被明确分类和记录。
- [x] KeyValues 风格 `.cfg` 有 fixture 和测试覆盖。
- [x] 更多 `.vdf` / `.acf` / Steam 文本样本被纳入测试。
- [x] GoDoc、CHANGELOG、中文文档和 release checklist 完成。
- [x] fuzz 和 malformed corpus 明显增强。
- [x] directive / condition token 策略清楚且有测试。
- [ ] 核心包仍保持小 API、无 Steam 业务绑定、无危险自动文件读取。
