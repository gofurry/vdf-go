# Changelog

All notable changes to `vdf-go` are documented here.

The project is pre-1.0. Public APIs may still evolve, but changes should remain
small, documented, and covered by tests.

## Unreleased

### Changed

- Moved implementation files behind `internal/vdf` while preserving the public
  module-root import path.
- Grouped fixtures into `testdata/valid`, `testdata/malformed`, and
  `testdata/unsupported`.
- Added architecture documentation for repository layout and package boundaries.

## v0.2.4 - Small AST Editing API

### Added

- `Document.Clone` and `Node.Clone`.
- `Document.Append` and `Node.Append`.
- `Document.SetFirst` and `Node.SetFirst`.
- `Document.RemoveFirst` / `Node.RemoveFirst`.
- `Document.RemoveAll` / `Node.RemoveAll`.

### Notes

- `SetFirst` replaces only the first matching key and preserves later duplicate
  keys.
- `RemoveAll` is explicit about removing every matching key.
- The AST remains slice-based and preserves order and duplicates.

## v0.2.3 - Fuzz and Malformed Corpus Hardening

### Added

- Extended malformed fixture corpus.
- Resource-limit regression tests for depth and token-size failures.
- Error-message regression tests that guard against leaking raw sensitive input.
- Fixture-based parse and marshal benchmarks.

### Validated

- `go test ./...`
- `go vet ./...`
- `staticcheck ./...`
- `go test -cover ./...`
- `go test -run=FuzzParse -fuzz=FuzzParse -fuzztime=1m`

## v0.2.2 - GoDoc, Chinese Docs, and Release Hygiene

### Added

- Package-level documentation for pkg.go.dev.
- Example tests for parsing, duplicate key querying, and marshaling.
- Chinese usage documentation and release notes.
- Release checklist for local validation, fuzzing, docs, and tagging.

## v0.2.1 - Text KeyValues Compatibility

### Added

- Condition token handling for tokens such as `[$WIN32]` after values or
  objects. Tokens are accepted and discarded; conditions are not evaluated.
- Regression tests for CRLF, UTF-8, empty keys, empty objects, unknown escapes,
  line comments, and parser error messages.

### Changed

- Only `#include` and `#base` are treated as special directives.
- Other `#...` keys are parsed as ordinary keys.

## v0.2.0 - Steam Text Fixture Compatibility

### Added

- Sanitized fixtures for `config.vdf`, `loginusers.vdf`, another
  `appmanifest_*.acf`, KeyValues-style `.cfg`, and command-style `.cfg`
  boundary documentation.
- Chinese compatibility matrix for `.vdf`, `.acf`, KeyValues-style `.cfg`, and
  unsupported Source / console command-style `.cfg`.

## v0.1.x - Text VDF Foundation

### Added

- Text VDF / KeyValues parser.
- `Document` and `Node` AST that preserves duplicate keys and order.
- Parse entrypoints: `Parse`, `ParseString`, `ParseReader`, and `ParseFile`.
- Query helpers: `First`, `All`, and `Path`.
- Constructors: `NewDocument`, `NewNode`, and `NewValue`.
- Marshal entrypoints: `Marshal`, `MarshalString`, and `Write`.
- Value helpers: `AsString`, `AsInt`, `AsUint64`, `AsFloat64`, and `AsBool`.
- Parser options for escapes, comments, bare tokens, depth, token size, node
  count, and directive preservation.
- Parser errors with line, column, and byte offset.
- Tests, fixtures, fuzz target, benchmarks, examples, README, CI, security, and
  contributing docs.
