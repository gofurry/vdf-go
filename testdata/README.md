# Test Fixtures

Fixtures in this directory are small, sanitized samples used to verify parser
compatibility. They should not contain real account identifiers, auth tokens,
machine-specific paths, or private install metadata.

## Supported by the VDF / KeyValues parser

- `valid/*.vdf`: text VDF / KeyValues samples.
- `valid/*.acf`: Steam appmanifest-style KeyValues samples.
- `valid/sample_keyvalues.cfg`: a KeyValues-style `.cfg` sample.

## Intentionally unsupported by the VDF / KeyValues parser

- `unsupported/source_commands.cfg`: Source / console command-style config.
  This format is line-oriented command text, not a VDF object tree.

## Malformed corpus

Malformed fixtures are intentionally invalid and should fail with ordinary
parser errors. They live under `malformed/` and cover:

- missing closing braces;
- unterminated quoted strings;
- unexpected closing braces;
- missing directive values;
- unterminated condition tokens;
- excessive nesting under a low `MaxDepth`;
- oversized tokens under a low `MaxTokenBytes`.
