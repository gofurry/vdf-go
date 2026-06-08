package vdf

// Node is a VDF / KeyValues entry.
//
// Children != nil means the node is an object. Children == nil means the node
// is a scalar value. This distinction preserves empty objects.
type Node struct {
	Key      string
	Value    string
	Children []*Node
}

// NewNode creates an object node with children.
func NewNode(key string, children ...*Node) *Node {
	if children == nil {
		children = []*Node{}
	}
	return &Node{Key: key, Children: children}
}

// NewValue creates a scalar value node.
func NewValue(key, value string) *Node {
	return &Node{Key: key, Value: value}
}

// IsObject reports whether n is an object node.
func (n *Node) IsObject() bool {
	return n != nil && n.Children != nil
}

// IsValue reports whether n is a scalar value node.
func (n *Node) IsValue() bool {
	return n != nil && n.Children == nil
}

// String returns the node value.
func (n *Node) String() string {
	return n.AsString()
}

// First returns the first child node with key.
func (n *Node) First(key string) *Node {
	if n == nil {
		return nil
	}
	return firstNode(n.Children, key)
}

// All returns all child nodes with key.
func (n *Node) All(key string) []*Node {
	if n == nil {
		return nil
	}
	return allNodes(n.Children, key)
}

// Path returns the first descendant at the provided key path.
func (n *Node) Path(keys ...string) *Node {
	if n == nil || len(keys) == 0 {
		return n
	}
	child := n.First(keys[0])
	if child == nil {
		return nil
	}
	return child.Path(keys[1:]...)
}

func firstNode(nodes []*Node, key string) *Node {
	for _, node := range nodes {
		if node != nil && node.Key == key {
			return node
		}
	}
	return nil
}

func allNodes(nodes []*Node, key string) []*Node {
	var matches []*Node
	for _, node := range nodes {
		if node != nil && node.Key == key {
			matches = append(matches, node)
		}
	}
	return matches
}
