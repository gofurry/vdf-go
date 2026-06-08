package vdf

import (
	"strings"
	"unicode/utf8"
)

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenString
	tokenBare
	tokenLBrace
	tokenRBrace
	tokenDirective
	tokenCondition
)

type position struct {
	line   int
	column int
	offset int
}

type token struct {
	kind tokenKind
	text string
	pos  position
}

type lexer struct {
	input []byte
	cfg   Config
	pos   int
	line  int
	col   int
}

func newLexer(input []byte, cfg Config) *lexer {
	return &lexer{input: input, cfg: cfg, line: 1, col: 1}
}

func (l *lexer) next() (token, error) {
	if err := l.skipSpaceAndComments(); err != nil {
		return token{}, err
	}
	pos := l.position()
	if l.eof() {
		return token{kind: tokenEOF, pos: pos}, nil
	}

	switch l.peekByte() {
	case '{':
		l.advanceByte()
		return token{kind: tokenLBrace, text: "{", pos: pos}, nil
	case '}':
		l.advanceByte()
		return token{kind: tokenRBrace, text: "}", pos: pos}, nil
	case '"':
		return l.scanQuoted()
	case '#':
		return l.scanDirective()
	case '[':
		return l.scanCondition()
	default:
		if !l.cfg.AllowBareTokens {
			return token{}, newParseError(pos, "unexpected bare token")
		}
		return l.scanBare()
	}
}

func (l *lexer) skipSpaceAndComments() error {
	for {
		for !l.eof() {
			switch l.peekByte() {
			case ' ', '\t', '\r', '\n':
				l.advanceByte()
			default:
				goto comments
			}
		}
		return nil

	comments:
		if l.cfg.AllowComments && l.hasPrefix("//") {
			for !l.eof() && l.peekByte() != '\n' {
				l.advanceByte()
			}
			continue
		}
		return nil
	}
}

func (l *lexer) scanQuoted() (token, error) {
	start := l.position()
	l.advanceByte()
	var b strings.Builder
	for !l.eof() {
		ch := l.peekByte()
		if ch == '"' {
			l.advanceByte()
			return token{kind: tokenString, text: b.String(), pos: start}, nil
		}
		if ch == '\\' && l.cfg.EscapeSequences {
			l.advanceByte()
			if l.eof() {
				return token{}, newParseError(start, "unterminated quoted string")
			}
			esc := l.peekByte()
			switch esc {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			case '"':
				b.WriteByte('"')
			default:
				b.WriteByte('\\')
				b.WriteByte(esc)
			}
			l.advanceByte()
		} else {
			r, size := utf8.DecodeRune(l.input[l.pos:])
			if r == utf8.RuneError && size == 1 {
				b.WriteByte(ch)
				l.advanceByte()
			} else {
				b.WriteRune(r)
				l.advanceBytes(size)
			}
		}
		if err := l.checkTokenSize(start, b.Len()); err != nil {
			return token{}, err
		}
	}
	return token{}, newParseError(start, "unterminated quoted string")
}

func (l *lexer) scanBare() (token, error) {
	start := l.position()
	begin := l.pos
	for !l.eof() {
		if l.startsComment() {
			break
		}
		switch l.peekByte() {
		case ' ', '\t', '\r', '\n', '{', '}':
			goto done
		default:
			l.advanceByte()
			if err := l.checkTokenSize(start, l.pos-begin); err != nil {
				return token{}, err
			}
		}
	}
done:
	if l.pos == begin {
		return token{}, newParseError(start, "unexpected character %q", l.peekByte())
	}
	return token{kind: tokenBare, text: string(l.input[begin:l.pos]), pos: start}, nil
}

func (l *lexer) scanDirective() (token, error) {
	start := l.position()
	begin := l.pos
	l.advanceByte()
	for !l.eof() {
		switch l.peekByte() {
		case ' ', '\t', '\r', '\n', '{', '}':
			goto done
		default:
			l.advanceByte()
			if err := l.checkTokenSize(start, l.pos-begin); err != nil {
				return token{}, err
			}
		}
	}
done:
	return token{kind: tokenDirective, text: string(l.input[begin:l.pos]), pos: start}, nil
}

func (l *lexer) scanCondition() (token, error) {
	start := l.position()
	begin := l.pos
	l.advanceByte()
	for !l.eof() {
		if l.peekByte() == ']' {
			l.advanceByte()
			return token{kind: tokenCondition, text: string(l.input[begin:l.pos]), pos: start}, nil
		}
		switch l.peekByte() {
		case '\r', '\n':
			return token{}, newParseError(start, "unterminated condition token")
		default:
			l.advanceByte()
			if err := l.checkTokenSize(start, l.pos-begin); err != nil {
				return token{}, err
			}
		}
	}
	return token{}, newParseError(start, "unterminated condition token")
}

func (l *lexer) checkTokenSize(pos position, size int) error {
	if l.cfg.MaxTokenBytes > 0 && size > l.cfg.MaxTokenBytes {
		return newParseError(pos, "token exceeds maximum size")
	}
	return nil
}

func (l *lexer) startsComment() bool {
	return l.cfg.AllowComments && l.hasPrefix("//")
}

func (l *lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(string(l.input[l.pos:]), prefix)
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) peekByte() byte {
	if l.eof() {
		return 0
	}
	return l.input[l.pos]
}

func (l *lexer) advanceByte() {
	if l.eof() {
		return
	}
	b := l.input[l.pos]
	l.pos++
	if b == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
}

func (l *lexer) advanceBytes(n int) {
	for i := 0; i < n; i++ {
		l.advanceByte()
	}
}

func (l *lexer) position() position {
	return position{line: l.line, column: l.col, offset: l.pos}
}
