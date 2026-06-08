package vdf

import "fmt"

// ParseError describes a parser failure with source position.
type ParseError struct {
	Message string
	Line    int
	Column  int
	Offset  int64
}

// Error returns a human-readable parser error.
func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("vdf: %s at line %d, column %d", e.Message, e.Line, e.Column)
	}
	return "vdf: " + e.Message
}

func newParseError(pos position, format string, args ...any) *ParseError {
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    pos.line,
		Column:  pos.column,
		Offset:  int64(pos.offset),
	}
}
