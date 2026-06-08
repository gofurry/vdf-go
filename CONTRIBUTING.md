# Contributing

Thanks for helping improve `vdf-go`.

## Development

Run the standard checks before sending changes:

```sh
gofmt -w *.go examples/*/*.go
go vet ./...
go test ./...
```

If `staticcheck` is available:

```sh
staticcheck ./...
```

## Design Principles

- Keep the public API small.
- Preserve node order and duplicate keys.
- Do not add Steam business logic to the core package.
- Do not automatically read files referenced by `#include` or `#base`.
- Prefer standard library dependencies.
- Add tests for public behavior.

## Tests

Fixtures live in `testdata/`. Add focused fixtures for real VDF shapes and
malformed inputs when changing parser behavior.

Fuzz tests should keep parser failures as ordinary errors and must not require
external services.
