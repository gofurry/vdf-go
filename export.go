package vdf

import (
	"io"

	core "github.com/gofurry/vdf-go/internal/vdf"
)

type Document = core.Document
type Node = core.Node

type Option = core.Option
type Config = core.Config
type ParseError = core.ParseError

type EncodeOption = core.EncodeOption
type EncodeConfig = core.EncodeConfig

func NewDocument(nodes ...*Node) *Document {
	return core.NewDocument(nodes...)
}

func NewNode(key string, children ...*Node) *Node {
	return core.NewNode(key, children...)
}

func NewValue(key, value string) *Node {
	return core.NewValue(key, value)
}

func DefaultConfig() Config {
	return core.DefaultConfig()
}

func WithEscapeSequences(enabled bool) Option {
	return core.WithEscapeSequences(enabled)
}

func WithMaxDepth(max int) Option {
	return core.WithMaxDepth(max)
}

func WithMaxTokenBytes(max int) Option {
	return core.WithMaxTokenBytes(max)
}

func WithMaxNodes(max int) Option {
	return core.WithMaxNodes(max)
}

func WithBareTokens(enabled bool) Option {
	return core.WithBareTokens(enabled)
}

func WithComments(enabled bool) Option {
	return core.WithComments(enabled)
}

func WithPreserveDirectives(enabled bool) Option {
	return core.WithPreserveDirectives(enabled)
}

func Parse(data []byte, opts ...Option) (*Document, error) {
	return core.Parse(data, opts...)
}

func ParseString(s string, opts ...Option) (*Document, error) {
	return core.ParseString(s, opts...)
}

func ParseReader(r io.Reader, opts ...Option) (*Document, error) {
	return core.ParseReader(r, opts...)
}

func ParseReaderLimit(r io.Reader, maxBytes int64, opts ...Option) (*Document, error) {
	return core.ParseReaderLimit(r, maxBytes, opts...)
}

func ParseFile(path string, opts ...Option) (*Document, error) {
	return core.ParseFile(path, opts...)
}

func DefaultEncodeConfig() EncodeConfig {
	return core.DefaultEncodeConfig()
}

func WithIndent(indent string) EncodeOption {
	return core.WithIndent(indent)
}

func WithQuoteKeys(enabled bool) EncodeOption {
	return core.WithQuoteKeys(enabled)
}

func WithQuoteValues(enabled bool) EncodeOption {
	return core.WithQuoteValues(enabled)
}

func WithSortKeys(enabled bool) EncodeOption {
	return core.WithSortKeys(enabled)
}

func Marshal(doc *Document, opts ...EncodeOption) ([]byte, error) {
	return core.Marshal(doc, opts...)
}

func MarshalString(doc *Document, opts ...EncodeOption) (string, error) {
	return core.MarshalString(doc, opts...)
}

func Write(w io.Writer, doc *Document, opts ...EncodeOption) error {
	return core.Write(w, doc, opts...)
}
