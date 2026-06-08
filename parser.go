package vdf

import (
	"fmt"
	"io"
	"os"
)

// Parse parses VDF / KeyValues data.
func Parse(data []byte, opts ...Option) (*Document, error) {
	cfg := applyOptions(opts)
	p := newParser(data, cfg)
	return p.parseDocument()
}

// ParseString parses VDF / KeyValues text.
func ParseString(s string, opts ...Option) (*Document, error) {
	return Parse([]byte(s), opts...)
}

// ParseReader reads all data from r and parses it.
func ParseReader(r io.Reader, opts ...Option) (*Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("vdf: read input: %w", err)
	}
	return Parse(data, opts...)
}

// ParseFile reads path and parses it. It does not expand #include or #base.
func ParseFile(path string, opts ...Option) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("vdf: read %q: %w", path, err)
	}
	return Parse(data, opts...)
}

type parser struct {
	lex       *lexer
	cfg       Config
	lookahead token
	has       bool
	nodes     int
}

func newParser(input []byte, cfg Config) *parser {
	return &parser{lex: newLexer(input, cfg), cfg: cfg}
}

func (p *parser) parseDocument() (*Document, error) {
	var nodes []*Node
	for {
		tok, err := p.peek()
		if err != nil {
			return nil, err
		}
		if tok.kind == tokenEOF {
			return &Document{Nodes: nodes}, nil
		}
		if tok.kind == tokenRBrace {
			return nil, newParseError(tok.pos, "unexpected %q", tok.text)
		}
		node, keep, err := p.parseNode(0)
		if err != nil {
			return nil, err
		}
		if keep {
			nodes = append(nodes, node)
		}
	}
}

func (p *parser) parseChildren(depth int) ([]*Node, error) {
	if p.cfg.MaxDepth > 0 && depth > p.cfg.MaxDepth {
		tok, _ := p.peek()
		return nil, newParseError(tok.pos, "maximum depth exceeded")
	}
	nodes := []*Node{}
	for {
		tok, err := p.peek()
		if err != nil {
			return nil, err
		}
		switch tok.kind {
		case tokenEOF:
			return nil, newParseError(tok.pos, "missing closing brace")
		case tokenRBrace:
			_, _ = p.next()
			return nodes, nil
		default:
			node, keep, err := p.parseNode(depth)
			if err != nil {
				return nil, err
			}
			if keep {
				nodes = append(nodes, node)
			}
		}
	}
}

func (p *parser) parseNode(depth int) (*Node, bool, error) {
	key, err := p.next()
	if err != nil {
		return nil, false, err
	}
	if !isKeyToken(key.kind) {
		return nil, false, newParseError(key.pos, "expected key")
	}

	value, err := p.next()
	if err != nil {
		return nil, false, err
	}
	switch value.kind {
	case tokenLBrace:
		if err := p.countNode(key.pos); err != nil {
			return nil, false, err
		}
		children, err := p.parseChildren(depth + 1)
		if err != nil {
			return nil, false, err
		}
		return &Node{Key: key.text, Children: children}, true, nil
	case tokenString, tokenBare:
		if key.kind == tokenDirective && !p.cfg.PreserveDirectives {
			return nil, false, nil
		}
		if err := p.countNode(key.pos); err != nil {
			return nil, false, err
		}
		return &Node{Key: key.text, Value: value.text}, true, nil
	case tokenEOF:
		return nil, false, newParseError(key.pos, "expected value or %q after key %q", "{", key.text)
	default:
		return nil, false, newParseError(value.pos, "expected value or %q after key %q", "{", key.text)
	}
}

func (p *parser) countNode(pos position) error {
	p.nodes++
	if p.cfg.MaxNodes > 0 && p.nodes > p.cfg.MaxNodes {
		return newParseError(pos, "node count exceeds maximum")
	}
	return nil
}

func (p *parser) peek() (token, error) {
	if p.has {
		return p.lookahead, nil
	}
	tok, err := p.lex.next()
	if err != nil {
		return token{}, err
	}
	p.lookahead = tok
	p.has = true
	return tok, nil
}

func (p *parser) next() (token, error) {
	if p.has {
		p.has = false
		return p.lookahead, nil
	}
	return p.lex.next()
}

func isKeyToken(kind tokenKind) bool {
	return kind == tokenString || kind == tokenBare || kind == tokenDirective
}
