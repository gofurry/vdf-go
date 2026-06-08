package vdf

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// EncodeOption configures VDF encoding.
type EncodeOption func(*EncodeConfig)

// EncodeConfig controls marshal output.
type EncodeConfig struct {
	Indent      string
	QuoteKeys   bool
	QuoteValues bool
	SortKeys    bool
}

// DefaultEncodeConfig returns encoder defaults.
func DefaultEncodeConfig() EncodeConfig {
	return EncodeConfig{
		Indent:      "\t",
		QuoteKeys:   true,
		QuoteValues: true,
		SortKeys:    false,
	}
}

// WithIndent sets one indentation unit.
func WithIndent(indent string) EncodeOption {
	return func(c *EncodeConfig) { c.Indent = indent }
}

// WithQuoteKeys controls whether keys are quoted.
func WithQuoteKeys(enabled bool) EncodeOption {
	return func(c *EncodeConfig) { c.QuoteKeys = enabled }
}

// WithQuoteValues controls whether scalar values are quoted.
func WithQuoteValues(enabled bool) EncodeOption {
	return func(c *EncodeConfig) { c.QuoteValues = enabled }
}

// WithSortKeys controls whether sibling nodes are sorted by key.
//
// Sorting is stable, so duplicate keys keep their relative order.
func WithSortKeys(enabled bool) EncodeOption {
	return func(c *EncodeConfig) { c.SortKeys = enabled }
}

// Marshal encodes doc as stable, readable VDF.
func Marshal(doc *Document, opts ...EncodeOption) ([]byte, error) {
	var buf bytes.Buffer
	if err := Write(&buf, doc, opts...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalString encodes doc as a string.
func MarshalString(doc *Document, opts ...EncodeOption) (string, error) {
	data, err := Marshal(doc, opts...)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Write encodes doc to w.
func Write(w io.Writer, doc *Document, opts ...EncodeOption) error {
	if w == nil {
		return fmt.Errorf("vdf: nil writer")
	}
	if doc == nil {
		return fmt.Errorf("vdf: nil document")
	}
	cfg := applyEncodeOptions(opts)
	enc := encoder{w: w, cfg: cfg}
	for _, node := range sortedNodes(doc.Nodes, cfg.SortKeys) {
		if node == nil {
			continue
		}
		if err := enc.writeNode(node, 0); err != nil {
			return err
		}
	}
	return nil
}

type encoder struct {
	w   io.Writer
	cfg EncodeConfig
}

func (e *encoder) writeNode(node *Node, depth int) error {
	if _, err := io.WriteString(e.w, strings.Repeat(e.cfg.Indent, depth)); err != nil {
		return err
	}
	if _, err := io.WriteString(e.w, e.formatToken(node.Key, e.cfg.QuoteKeys)); err != nil {
		return err
	}
	if node.IsObject() {
		if _, err := io.WriteString(e.w, "\n"+strings.Repeat(e.cfg.Indent, depth)+"{\n"); err != nil {
			return err
		}
		for _, child := range sortedNodes(node.Children, e.cfg.SortKeys) {
			if child == nil {
				continue
			}
			if err := e.writeNode(child, depth+1); err != nil {
				return err
			}
		}
		_, err := io.WriteString(e.w, strings.Repeat(e.cfg.Indent, depth)+"}\n")
		return err
	}
	if _, err := io.WriteString(e.w, "\t"); err != nil {
		return err
	}
	if _, err := io.WriteString(e.w, e.formatToken(node.Value, e.cfg.QuoteValues)); err != nil {
		return err
	}
	_, err := io.WriteString(e.w, "\n")
	return err
}

func (e *encoder) formatToken(s string, quote bool) string {
	if !quote {
		return s
	}
	return `"` + escapeString(s) + `"`
}

func escapeString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func applyEncodeOptions(opts []EncodeOption) EncodeConfig {
	cfg := DefaultEncodeConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

func sortedNodes(nodes []*Node, sortKeys bool) []*Node {
	if !sortKeys || len(nodes) < 2 {
		return nodes
	}
	out := append([]*Node(nil), nodes...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i] == nil {
			return false
		}
		if out[j] == nil {
			return true
		}
		return out[i].Key < out[j].Key
	})
	return out
}
