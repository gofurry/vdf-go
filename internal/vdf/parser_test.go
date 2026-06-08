package vdf

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestParseBasicFeatures(t *testing.T) {
	input := `"AppState"
{
	"appid"	"730"
	name "Counter-Strike 2"
	"installdir" "Counter-Strike Global Offensive"
	"windows_path" "C:\\Program Files (x86)\\Steam"
	// comment
	"escaped" "line\nquote\"tab\tbackslash\\"
}
"root2" "value2"
`
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	app := doc.First("AppState")
	if app == nil || !app.IsObject() {
		t.Fatalf("AppState object not found")
	}
	if got := app.First("appid").Value; got != "730" {
		t.Fatalf("appid = %q", got)
	}
	if got := app.First("windows_path").Value; got != `C:\Program Files (x86)\Steam` {
		t.Fatalf("windows_path = %q", got)
	}
	if got := app.First("escaped").Value; got != "line\nquote\"tab\tbackslash\\" {
		t.Fatalf("escaped = %q", got)
	}
	if got := doc.First("root2").Value; got != "value2" {
		t.Fatalf("root2 = %q", got)
	}
}

func TestParseDuplicateKeysAndPath(t *testing.T) {
	doc, err := ParseString(`"root" { "item" "one" "item" "two" "child" { "leaf" "ok" } }`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	root := doc.First("root")
	items := root.All("item")
	if len(items) != 2 {
		t.Fatalf("duplicate item count = %d", len(items))
	}
	if items[0].Value != "one" || items[1].Value != "two" {
		t.Fatalf("items out of order: %#v", items)
	}
	if got := doc.Path("root", "child", "leaf").Value; got != "ok" {
		t.Fatalf("path value = %q", got)
	}
}

func TestParseEmptyValueAndEmptyObject(t *testing.T) {
	doc, err := ParseString(`"root" { "empty_value" "" "empty_object" { } }`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	root := doc.First("root")
	if !root.First("empty_value").IsValue() {
		t.Fatalf("empty_value should be value")
	}
	if !root.First("empty_object").IsObject() {
		t.Fatalf("empty_object should be object")
	}
}

func TestDirectivesDefaultIgnoredAndPreserved(t *testing.T) {
	input := `#include "base.vdf" "root" { #base "defaults.vdf" "k" "v" } #custom "kept"`
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if doc.First("#include") != nil {
		t.Fatalf("directive should be ignored by default")
	}
	if doc.First("root").First("#base") != nil {
		t.Fatalf("nested directive should be ignored by default")
	}
	if got := doc.First("#custom").Value; got != "kept" {
		t.Fatalf("unknown directive-like key should be kept, got %q", got)
	}

	doc, err = ParseString(input, WithPreserveDirectives(true))
	if err != nil {
		t.Fatalf("ParseString(preserve) error = %v", err)
	}
	if got := doc.First("#include").Value; got != "base.vdf" {
		t.Fatalf("include = %q", got)
	}
	if got := doc.First("root").First("#base").Value; got != "defaults.vdf" {
		t.Fatalf("base = %q", got)
	}
}

func TestConditionTokensAreAcceptedButNotEvaluated(t *testing.T) {
	input := `"root" { "windows" "kept" [$WIN32] "linux" "also kept" [$LINUX] } [$X360] "next" "ok"`
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	root := doc.First("root")
	if got := root.First("windows").Value; got != "kept" {
		t.Fatalf("windows = %q", got)
	}
	if got := root.First("linux").Value; got != "also kept" {
		t.Fatalf("linux = %q", got)
	}
	if got := doc.First("next").Value; got != "ok" {
		t.Fatalf("next = %q", got)
	}
}

func TestUnexpectedConditionErrors(t *testing.T) {
	tests := []string{
		`[$WIN32] "key" "value"`,
		`"root" { [$WIN32] "key" "value" }`,
		`"key" "value" [$WIN32`,
	}
	for _, input := range tests {
		if _, err := ParseString(input); err == nil {
			t.Fatalf("ParseString(%q) succeeded", input)
		}
	}
}

func TestTextKeyValuesEdgeCases(t *testing.T) {
	input := "\"\" \"empty key\"\r\n\"utf8\" \"简体中文\"\r\n\"empty_object\" { }\r\n\"unknown_escape\" \"a\\zb\" // line comment\r\n\"next\" \"ok\""
	doc, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if got := doc.First("").Value; got != "empty key" {
		t.Fatalf("empty key value = %q", got)
	}
	if got := doc.First("utf8").Value; got != "简体中文" {
		t.Fatalf("utf8 value = %q", got)
	}
	if !doc.First("empty_object").IsObject() {
		t.Fatalf("empty_object should be object")
	}
	if got := doc.First("unknown_escape").Value; got != `a\zb` {
		t.Fatalf("unknown_escape = %q", got)
	}
	if got := doc.First("next").Value; got != "ok" {
		t.Fatalf("next = %q", got)
	}
}

func TestParseErrorMessagesRegression(t *testing.T) {
	tests := []struct {
		input string
		want  string
		line  int
	}{
		{input: "\"root\"\r\n{\r\n\"k\" \"v\"", want: "missing closing brace", line: 3},
		{input: `"root" { "k" "unterminated`, want: "unterminated quoted string", line: 1},
		{input: `"root" { "k" }`, want: `expected value or "{" after key "k"`, line: 1},
	}
	for _, tt := range tests {
		_, err := ParseString(tt.input)
		if err == nil {
			t.Fatalf("ParseString(%q) succeeded", tt.input)
		}
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("error %T is not *ParseError", err)
		}
		if !strings.Contains(parseErr.Message, tt.want) {
			t.Fatalf("message = %q, want contains %q", parseErr.Message, tt.want)
		}
		if parseErr.Line != tt.line {
			t.Fatalf("line = %d, want %d", parseErr.Line, tt.line)
		}
	}
}

func TestParseOptions(t *testing.T) {
	if _, err := ParseString(`key value`, WithBareTokens(false)); err == nil {
		t.Fatalf("expected bare token error")
	}
	doc, err := ParseString(`"k" "v" // comment`, WithComments(false))
	if err != nil {
		t.Fatalf("ParseString(comments disabled) error = %v", err)
	}
	if got := doc.First("//").Value; got != "comment" {
		t.Fatalf("disabled comment token = %q", got)
	}
	doc, err = ParseString(`"k" "a\nb"`, WithEscapeSequences(false))
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	if got := doc.First("k").Value; got != `a\nb` {
		t.Fatalf("escape-disabled value = %q", got)
	}
}

func TestParseErrorsHavePosition(t *testing.T) {
	tests := []string{
		`"root" { "k" "v"`,
		`"root" { "k" "v" } }`,
		`"root" { "k" "unterminated`,
		`"root" { "k" }`,
	}
	for _, input := range tests {
		_, err := ParseString(input)
		if err == nil {
			t.Fatalf("ParseString(%q) succeeded", input)
		}
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("error %T is not *ParseError", err)
		}
		if parseErr.Line <= 0 || parseErr.Column <= 0 || parseErr.Offset < 0 {
			t.Fatalf("bad position: %#v", parseErr)
		}
	}
}

func TestParseLimits(t *testing.T) {
	if _, err := ParseString(`"a" { "b" { "c" "d" } }`, WithMaxDepth(1)); err == nil {
		t.Fatalf("expected max depth error")
	}
	if _, err := ParseString(`"long" "value"`, WithMaxTokenBytes(3)); err == nil {
		t.Fatalf("expected max token error")
	}
	if _, err := ParseString(`"a" "1" "b" "2"`, WithMaxNodes(1)); err == nil {
		t.Fatalf("expected max nodes error")
	}
}

func TestParseReader(t *testing.T) {
	doc, err := ParseReader(strings.NewReader(`"k" "v"`))
	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}
	if got := doc.First("k").Value; got != "v" {
		t.Fatalf("value = %q", got)
	}
}

func TestParseReaderLimit(t *testing.T) {
	doc, err := ParseReaderLimit(strings.NewReader(`"k" "v"`), 7)
	if err != nil {
		t.Fatalf("ParseReaderLimit() error = %v", err)
	}
	if got := doc.First("k").Value; got != "v" {
		t.Fatalf("value = %q", got)
	}

	if _, err := ParseReaderLimit(strings.NewReader(`"k" "v"`), 6); err == nil {
		t.Fatalf("expected size limit error")
	}
	if _, err := ParseReaderLimit(strings.NewReader(`"k" "v"`), -1); err == nil {
		t.Fatalf("expected negative limit error")
	}
	if _, err := ParseReaderLimit(strings.NewReader(``), 0); err != nil {
		t.Fatalf("empty input with zero limit should parse: %v", err)
	}
}

func TestFixtures(t *testing.T) {
	tests := []struct {
		path string
		key  string
	}{
		{"../../testdata/valid/source-keyvalues/simple.vdf", "root"},
		{"../../testdata/valid/source-keyvalues/nested.vdf", "root"},
		{"../../testdata/valid/source-keyvalues/duplicate_keys.vdf", "root"},
		{"../../testdata/valid/source-keyvalues/keyvalues_cfg.cfg", "SampleConfig"},
		{"../../testdata/valid/source-keyvalues/gameinfo.txt", "GameInfo"},
		{"../../testdata/valid/source-keyvalues/material_basic.vmt", "LightmappedGeneric"},
		{"../../testdata/valid/source-keyvalues/resource_menu.res", "Resource/UI/MainMenu.res"},
		{"../../testdata/valid/source-keyvalues/scripts_sounds.txt", "Weapon.SampleFire"},
		{"../../testdata/valid/steam-client/libraryfolders_fixture.vdf", "libraryfolders"},
		{"../../testdata/valid/steam-client/config_fixture.vdf", "InstallConfigStore"},
		{"../../testdata/valid/steam-client/loginusers_fixture.vdf", "users"},
		{"../../testdata/valid/steam-client/libraryfolders_steamapps_sanitized.vdf", "libraryfolders"},
		{"../../testdata/valid/steam-client/libraryfolders_config_sanitized.vdf", "libraryfolders"},
		{"../../testdata/valid/steam-client/config_sanitized.vdf", "InstallConfigStore"},
		{"../../testdata/valid/steam-client/dialog_config_sanitized.vdf", "UserConfigData"},
		{"../../testdata/valid/steam-client/loginusers_sanitized.vdf", "users"},
		{"../../testdata/valid/appmanifest/appmanifest_730.acf", "AppState"},
		{"../../testdata/valid/appmanifest/appmanifest_570.acf", "AppState"},
		{"../../testdata/valid/directives/include_ignored.vdf", "root"},
		{"../../testdata/valid/directives/base_ignored.vdf", "root"},
		{"../../testdata/valid/directives/preserve_directives.vdf", "root"},
		{"../../testdata/valid/conditions/win32_condition.vdf", "root"},
		{"../../testdata/valid/conditions/mixed_conditions.vdf", "root"},
	}
	for _, tt := range tests {
		doc, err := ParseFile(tt.path)
		if err != nil {
			t.Fatalf("ParseFile(%q) error = %v", tt.path, err)
		}
		if doc.First(tt.key) == nil {
			t.Fatalf("ParseFile(%q) missing %q", tt.path, tt.key)
		}
	}

	doc, err := ParseFile("../../testdata/valid/source-keyvalues/duplicate_keys.vdf")
	if err != nil {
		t.Fatalf("ParseFile(duplicate) error = %v", err)
	}
	if got := len(doc.First("root").All("item")); got != 3 {
		t.Fatalf("fixture duplicate count = %d", got)
	}

	app, err := ParseFile("../../testdata/valid/appmanifest/appmanifest_730.acf")
	if err != nil {
		t.Fatalf("ParseFile(appmanifest) error = %v", err)
	}
	if got := app.Path("AppState", "InstalledDepots", "731", "manifest").Value; got == "" {
		t.Fatalf("missing depot manifest")
	}

	cfg, err := ParseFile("../../testdata/valid/source-keyvalues/keyvalues_cfg.cfg")
	if err != nil {
		t.Fatalf("ParseFile(sample_keyvalues.cfg) error = %v", err)
	}
	if got := cfg.Path("SampleConfig", "Profile", "Mode").Value; got != "safe" {
		t.Fatalf("cfg mode = %q", got)
	}
}

func TestSanitizedSteamClientFixtures(t *testing.T) {
	config, err := ParseFile("../../testdata/valid/steam-client/config_sanitized.vdf")
	if err != nil {
		t.Fatalf("ParseFile(config_sanitized.vdf) error = %v", err)
	}
	if got := config.Path("InstallConfigStore", "Software", "Valve", "Steam", "cip").Value; got != "203.0.113.10" {
		t.Fatalf("sanitized cip = %q", got)
	}
	if got := config.Path("InstallConfigStore", "Software", "Valve", "Steam", "Accounts", "sample_account", "SteamID").Value; got != "76561198000000000" {
		t.Fatalf("sanitized account SteamID = %q", got)
	}
	if node := config.Path("InstallConfigStore", "Software", "Valve", "Steam", "CMWebSocket", "cm1.example.steamserver.net:443"); node == nil || !node.IsObject() {
		t.Fatalf("missing sanitized CMWebSocket entry")
	}

	login, err := ParseFile("../../testdata/valid/steam-client/loginusers_sanitized.vdf")
	if err != nil {
		t.Fatalf("ParseFile(loginusers_sanitized.vdf) error = %v", err)
	}
	users := login.First("users")
	if users == nil || len(users.Children) != 2 {
		t.Fatalf("sanitized login users count = %d", len(users.Children))
	}
	if got := login.Path("users", "76561198000000000", "AccountName").Value; got != "sample_account" {
		t.Fatalf("sanitized login AccountName = %q", got)
	}
	if got := login.Path("users", "76561198000000001", "PersonaName").Value; got != "Sample Alt" {
		t.Fatalf("sanitized login PersonaName = %q", got)
	}

	library, err := ParseFile("../../testdata/valid/steam-client/libraryfolders_steamapps_sanitized.vdf")
	if err != nil {
		t.Fatalf("ParseFile(libraryfolders_steamapps_sanitized.vdf) error = %v", err)
	}
	if got := library.Path("libraryfolders", "0", "path").Value; got != `C:\Program Files (x86)\Steam` {
		t.Fatalf("sanitized library path = %q", got)
	}
	if !library.Path("libraryfolders", "0", "apps").IsObject() {
		t.Fatalf("sanitized library apps should be object")
	}
}

func TestDirectiveFixtures(t *testing.T) {
	doc, err := ParseFile("../../testdata/valid/directives/preserve_directives.vdf")
	if err != nil {
		t.Fatalf("ParseFile(preserve_directives.vdf) error = %v", err)
	}
	if doc.First("#include") != nil {
		t.Fatalf("directive should be ignored by default")
	}

	doc, err = ParseFile("../../testdata/valid/directives/preserve_directives.vdf", WithPreserveDirectives(true))
	if err != nil {
		t.Fatalf("ParseFile(preserve_directives.vdf preserve) error = %v", err)
	}
	if got := doc.First("#include").Value; got != "shared_defaults.vdf" {
		t.Fatalf("include directive = %q", got)
	}
	if got := doc.First("root").First("#base").Value; got != "base_defaults.vdf" {
		t.Fatalf("base directive = %q", got)
	}
}

func TestCommandStyleCFGIsNotVDFScope(t *testing.T) {
	data, err := os.ReadFile("../../testdata/unsupported/source_commands.cfg")
	if err != nil {
		t.Fatalf("ReadFile(source_commands.cfg) error = %v", err)
	}
	text := string(data)
	for _, want := range []string{"bind", "con_enable", "not a VDF / KeyValues object tree"} {
		if !strings.Contains(text, want) {
			t.Fatalf("source_commands.cfg missing %q", want)
		}
	}
}

func TestMalformedFixtures(t *testing.T) {
	for _, path := range []string{
		"../../testdata/malformed/malformed_missing_brace.vdf",
		"../../testdata/malformed/malformed_unterminated_quote.vdf",
		"../../testdata/malformed/malformed_unexpected_closing_brace.vdf",
		"../../testdata/malformed/malformed_directive_missing_value.vdf",
		"../../testdata/malformed/malformed_unterminated_condition.vdf",
	} {
		_, err := ParseFile(path)
		if err == nil {
			t.Fatalf("ParseFile(%q) succeeded", path)
		}
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("ParseFile(%q) error %T is not *ParseError", path, err)
		}
	}
}

func TestMalformedResourceLimitFixtures(t *testing.T) {
	tests := []struct {
		path string
		opts []Option
		want string
	}{
		{
			path: "../../testdata/malformed/malformed_deep_nesting.vdf",
			opts: []Option{WithMaxDepth(3)},
			want: "maximum depth exceeded",
		},
		{
			path: "../../testdata/malformed/malformed_token_too_large.vdf",
			opts: []Option{WithMaxTokenBytes(16)},
			want: "token exceeds maximum size",
		},
	}
	for _, tt := range tests {
		_, err := ParseFile(tt.path, tt.opts...)
		if err == nil {
			t.Fatalf("ParseFile(%q) succeeded", tt.path)
		}
		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("ParseFile(%q) error %T is not *ParseError", tt.path, err)
		}
		if !strings.Contains(parseErr.Message, tt.want) {
			t.Fatalf("ParseFile(%q) message = %q, want %q", tt.path, parseErr.Message, tt.want)
		}
	}
}

func TestParseErrorsDoNotEchoSensitiveInput(t *testing.T) {
	tests := []struct {
		input  string
		opts   []Option
		secret string
	}{
		{
			input:  `"root" { "token" "super-secret-value`,
			secret: "super-secret-value",
		},
		{
			input:  `"super-secret-token-value" "x"`,
			opts:   []Option{WithMaxTokenBytes(8)},
			secret: "super-secret-token-value",
		},
	}
	for _, tt := range tests {
		_, err := ParseString(tt.input, tt.opts...)
		if err == nil {
			t.Fatalf("ParseString(%q) succeeded", tt.input)
		}
		if strings.Contains(err.Error(), tt.secret) {
			t.Fatalf("error leaked sensitive input %q: %v", tt.secret, err)
		}
	}
}

func TestParseFileReadError(t *testing.T) {
	path := "../../testdata/does-not-exist.vdf"
	if _, err := ParseFile(path); err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ParseFile(%q) error = %v", path, err)
	}
}
