# Architecture

`vdf-go` keeps the public import path stable while allowing the implementation
to evolve behind an internal boundary.

## Layout

```text
.
├── export.go              # public facade for github.com/gofurry/vdf-go
├── doc.go                 # package-level GoDoc
├── examples_test.go       # public facade examples
├── internal/
│   └── vdf/               # parser, encoder, AST, options, errors, editing, tests
├── testdata/
│   ├── valid/             # supported Steam / Source text KeyValues fixtures
│   │   ├── steam-client/
│   │   ├── appmanifest/
│   │   ├── source-keyvalues/
│   │   ├── directives/
│   │   └── conditions/
│   ├── malformed/         # invalid inputs for parser hardening
│   └── unsupported/       # documented non-goals and boundary samples
├── examples/              # runnable examples
└── docs/                  # English and Chinese documentation
```

## Public API Boundary

Users should import:

```go
import vdf "github.com/gofurry/vdf-go"
```

The root package re-exports the stable public API from `internal/vdf`. The
internal package is not part of the public contract and may be reorganized
without changing the user-facing import path.

## Why Not Move the Public Package?

Moving the public package into a subdirectory such as `vdf/` would change the
import path to `github.com/gofurry/vdf-go/vdf`. This repository intentionally
keeps the shorter module-root import path.

## Test Data Policy

Fixtures should stay small, sanitized, and categorized:

- `testdata/valid/steam-client`: sanitized real-world Steam client shapes;
- `testdata/valid/appmanifest`: Steam appmanifest `.acf` shapes;
- `testdata/valid/source-keyvalues`: Source / Valve KeyValues-style text files;
- `testdata/valid/directives`: `#include` / `#base` behavior samples;
- `testdata/valid/conditions`: condition-token behavior samples;
- `testdata/malformed`: invalid inputs that must fail safely;
- `testdata/unsupported`: examples that document format boundaries.

Do not commit real account identifiers, auth tokens, private installation paths,
or machine-specific metadata.

## Test Placement

Implementation tests live beside the implementation in `internal/vdf`, so
coverage reflects the parser, encoder, AST, and editing code directly. The root
package keeps only facade-level examples to verify the public import path.
