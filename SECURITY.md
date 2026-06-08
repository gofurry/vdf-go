# Security Policy

## Supported Versions

Security fixes are expected for the latest `v0.1.x` release while the project is
pre-1.0.

## Parser Safety Boundary

`vdf-go` parses text VDF / KeyValues data only. It does not execute directives,
does not expand `#include` or `#base`, and does not read files other than the
explicit path passed to `ParseFile`.

The parser includes configurable limits for nesting depth, token size, and node
count. Callers that accept untrusted input should keep these limits enabled.

## Reporting Issues

Please report security-sensitive issues privately to the repository owner before
opening public issues.

Include:

- affected version or commit;
- minimal input that demonstrates the issue;
- expected impact;
- whether the issue can be triggered by untrusted VDF input.
