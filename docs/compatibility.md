# Compatibility

This document describes the compatibility boundary for `vdf-go`.

## Supported

- Text VDF / Valve KeyValues files.
- Quoted keys and values.
- Unquoted keys and values when `AllowBareTokens` is enabled.
- Nested `{}` objects.
- Multiple root nodes.
- Duplicate sibling keys.
- Child order preservation.
- `//` line comments when `AllowComments` is enabled.
- Basic escapes in quoted strings: `\n`, `\t`, `\\`, and `\"`.
- Common Steam text files such as `libraryfolders.vdf`, `config.vdf`,
  `loginusers.vdf`, and `appmanifest_*.acf`.
- KeyValues-style `.cfg` files.
- Condition tokens such as `[$WIN32]` after values or objects. They are accepted
  and discarded.
- Parser errors with line, column, and byte offset.
- Stable readable marshal output.

## Not Supported

- Source / console command-style `.cfg` files, such as line-oriented `bind`,
  `alias`, and cvar config files.
- Binary VDF.
- `shortcuts.vdf`.
- Struct decoder or encoder.
- Steam local library scanning.
- Strongly typed Steam appmanifest or libraryfolders models.
- Comment round-trip preservation.
- Original whitespace and formatting preservation.
- Automatic `#include` or `#base` file loading.
- Conditional token evaluation such as `[$WIN32]`.
- A formatter CLI.

## Directives

`#include` and `#base` are treated conservatively. The parser never opens files
referenced by directives.

Default behavior:

- `#include` and `#base` are consumed and ignored;
- other `#...` keys are parsed as ordinary keys;
- no filesystem reads are performed;
- no path normalization is performed.

With `WithPreserveDirectives(true)`:

- directives are kept as ordinary value nodes;
- for `#include "base.vdf"`, the node key is `#include` and value is
  `base.vdf`;
- the referenced file is still not read.

## Condition Tokens

Condition tokens such as `[$WIN32]` are accepted after a value or object and then
discarded. The parser does not evaluate the condition, so the associated node is
parsed unconditionally.

## Ordering and Duplicate Keys

`vdf-go` intentionally does not model objects as maps. Sibling order and
duplicate keys are part of the parsed document.

Use `First` for convenience when the first matching key is enough. Use `All`
when duplicate keys matter.

## Formatting

Marshal output is designed to be stable and readable. It does not preserve
original comments, blank lines, spacing, or quote style.

## Chinese Documentation

The detailed Chinese compatibility matrix is maintained at
[`docs/zh/compatibility.md`](zh/compatibility.md).
