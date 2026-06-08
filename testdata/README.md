# Test Fixtures

Fixtures in this directory are small, sanitized samples used to verify parser
compatibility. They should not contain real account identifiers, auth tokens,
machine-specific paths, or private install metadata.

## Supported by the VDF / KeyValues parser

- `*.vdf`: text VDF / KeyValues samples.
- `*.acf`: Steam appmanifest-style KeyValues samples.
- `sample_keyvalues.cfg`: a KeyValues-style `.cfg` sample.

## Intentionally unsupported by the VDF / KeyValues parser

- `source_commands.cfg`: Source / console command-style config. This format is
  line-oriented command text, not a VDF object tree.

## Malformed corpus

Malformed fixtures are intentionally invalid and should fail with ordinary
parser errors:

- missing closing braces;
- unterminated quoted strings;
- unexpected closing braces;
- missing directive values;
- unterminated condition tokens;
- excessive nesting under a low `MaxDepth`;
- oversized tokens under a low `MaxTokenBytes`.
