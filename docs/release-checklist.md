# Release Checklist

Use this checklist before tagging a local release.

## 1. Working Tree

- [ ] Confirm the branch is correct.
- [ ] Confirm `git status --short` contains only intentional changes.
- [ ] Review `git diff` and `git diff --staged`.
- [ ] Check for private fixtures, real account IDs, auth tokens, credentials, or
      machine-specific paths.

## 2. Validation

Run:

```sh
gofmt -w .
go test ./...
go test -race ./...
go test -cover ./...
go vet ./...
staticcheck ./...
```

For parser-affecting releases, also run fuzzing:

```sh
go test -run=FuzzParse -fuzz=FuzzParse -fuzztime=1m
```

Promote useful fuzz discoveries into normal regression tests.

## 3. Documentation

- [ ] Update `README.md` when public behavior changes.
- [ ] Update `docs/compatibility.md` and `docs/zh/compatibility.md` when format
      support changes.
- [ ] Update `docs/zh/roadmap.md` when roadmap items are completed or deferred.
- [ ] Update `CHANGELOG.md`.
- [ ] Confirm examples still run.

## 4. Versioning

- [ ] Use `v0.x.y` before stable API freeze.
- [ ] Reserve `v1.0.0-alpha.x` for API-freeze candidates.
- [ ] Reserve `v1.0.0` for the first official stable release.

## 5. Tagging

After validation and commit:

```sh
git tag v0.x.y
git show v0.x.y
```

Push only when intentionally publishing:

```sh
git push origin dev
git push origin v0.x.y
```
