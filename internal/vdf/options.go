package vdf

// Option configures parsing.
type Option func(*Config)

// Config controls parser behavior and resource limits.
type Config struct {
	EscapeSequences    bool
	MaxDepth           int
	MaxTokenBytes      int
	MaxNodes           int
	AllowBareTokens    bool
	AllowComments      bool
	PreserveDirectives bool
}

// DefaultConfig returns parser defaults.
func DefaultConfig() Config {
	return Config{
		EscapeSequences:    true,
		MaxDepth:           128,
		MaxTokenBytes:      1 << 20,
		MaxNodes:           1_000_000,
		AllowBareTokens:    true,
		AllowComments:      true,
		PreserveDirectives: false,
	}
}

// WithEscapeSequences enables or disables quoted-string escape processing.
func WithEscapeSequences(enabled bool) Option {
	return func(c *Config) { c.EscapeSequences = enabled }
}

// WithMaxDepth sets the maximum object nesting depth. Values <= 0 disable the limit.
func WithMaxDepth(max int) Option {
	return func(c *Config) { c.MaxDepth = max }
}

// WithMaxTokenBytes sets the maximum bytes in one token. Values <= 0 disable the limit.
func WithMaxTokenBytes(max int) Option {
	return func(c *Config) { c.MaxTokenBytes = max }
}

// WithMaxNodes sets the maximum parsed node count. Values <= 0 disable the limit.
func WithMaxNodes(max int) Option {
	return func(c *Config) { c.MaxNodes = max }
}

// WithBareTokens enables or disables unquoted key/value tokens.
func WithBareTokens(enabled bool) Option {
	return func(c *Config) { c.AllowBareTokens = enabled }
}

// WithComments enables or disables // line comments.
func WithComments(enabled bool) Option {
	return func(c *Config) { c.AllowComments = enabled }
}

// WithPreserveDirectives controls whether #include/#base directives are kept as nodes.
//
// Directives are never executed and never read files.
func WithPreserveDirectives(enabled bool) Option {
	return func(c *Config) { c.PreserveDirectives = enabled }
}

func applyOptions(opts []Option) Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}
