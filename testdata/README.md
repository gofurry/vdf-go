# Test Fixtures

Fixtures in this directory are small, sanitized samples used to verify parser
compatibility. They should not contain real account identifiers, auth tokens,
machine-specific paths, or private install metadata.

## Supported by the VDF / KeyValues parser

- `valid/steam-client/`: sanitized Steam client text VDF samples, including
  `libraryfolders.vdf`, `config.vdf`, `DialogConfig.vdf`, and `loginusers.vdf`
  shapes.
- `valid/appmanifest/`: Steam appmanifest-style `.acf` samples.
- `valid/source-keyvalues/`: Source / Valve KeyValues-style text samples such
  as `.cfg`, `.txt`, `.vmt`, and `.res`.
- `valid/directives/`: `#include` / `#base` fixtures for default ignored mode
  and preserve mode.
- `valid/conditions/`: condition-token fixtures such as `[$WIN32]`.

## Intentionally unsupported by the VDF / KeyValues parser

- `unsupported/source_commands.cfg`: Source / console command-style config.
  This format is line-oriented command text, not a VDF object tree.
- `unsupported/shortcuts_binary_note.md`: documents why Steam `shortcuts.vdf`
  is outside text parser scope.
- `unsupported/binary_vdf_sample_note.md`: documents binary VDF as a future
  separate parser concern.

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
