package vdf

import (
	"fmt"
	"strconv"
	"strings"
)

// AsString returns the node's scalar value, or an empty string for nil nodes.
func (n *Node) AsString() string {
	if n == nil {
		return ""
	}
	return n.Value
}

// AsInt parses the node value as int.
func (n *Node) AsInt() (int, error) {
	v, err := strconv.Atoi(n.AsString())
	if err != nil {
		return 0, fmt.Errorf("vdf: parse %q as int: %w", n.AsString(), err)
	}
	return v, nil
}

// AsUint64 parses the node value as uint64.
func (n *Node) AsUint64() (uint64, error) {
	v, err := strconv.ParseUint(n.AsString(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("vdf: parse %q as uint64: %w", n.AsString(), err)
	}
	return v, nil
}

// AsFloat64 parses the node value as float64.
func (n *Node) AsFloat64() (float64, error) {
	v, err := strconv.ParseFloat(n.AsString(), 64)
	if err != nil {
		return 0, fmt.Errorf("vdf: parse %q as float64: %w", n.AsString(), err)
	}
	return v, nil
}

// AsBool parses common VDF convenience booleans.
//
// True values are 1, true, yes, and on. False values are 0, false, no, and off.
// This is a vdf-go convenience helper, not an official Valve type system.
func (n *Node) AsBool() (bool, error) {
	switch strings.ToLower(strings.TrimSpace(n.AsString())) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("vdf: parse %q as bool", n.AsString())
	}
}
