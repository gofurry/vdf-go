package vdf

// Clone returns a deep copy of d.
func (d *Document) Clone() *Document {
	if d == nil {
		return nil
	}
	return &Document{Nodes: cloneNodes(d.Nodes)}
}

// Append appends root nodes to d.
func (d *Document) Append(nodes ...*Node) {
	if d == nil {
		return
	}
	d.Nodes = append(d.Nodes, nodes...)
}

// SetFirst replaces the first root node with the same key as node.
//
// If no root node has the same key, node is appended. Duplicate keys after the
// first match are preserved. It returns true when an existing node was replaced.
func (d *Document) SetFirst(node *Node) bool {
	if d == nil || node == nil {
		return false
	}
	for i, existing := range d.Nodes {
		if existing != nil && existing.Key == node.Key {
			d.Nodes[i] = node
			return true
		}
	}
	d.Nodes = append(d.Nodes, node)
	return false
}

// RemoveFirst removes and returns the first root node with key.
func (d *Document) RemoveFirst(key string) *Node {
	if d == nil {
		return nil
	}
	removed, nodes := removeFirstNode(d.Nodes, key)
	d.Nodes = nodes
	return removed
}

// RemoveAll removes and returns all root nodes with key.
func (d *Document) RemoveAll(key string) []*Node {
	if d == nil {
		return nil
	}
	removed, nodes := removeAllNodes(d.Nodes, key)
	d.Nodes = nodes
	return removed
}

// Clone returns a deep copy of n.
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{
		Key:   n.Key,
		Value: n.Value,
	}
	if n.Children != nil {
		clone.Children = cloneNodes(n.Children)
	}
	return clone
}

// Append appends child nodes to n.
//
// If n is a value node, Append converts it to an object node and clears Value.
func (n *Node) Append(children ...*Node) {
	if n == nil {
		return
	}
	if n.Children == nil {
		n.Value = ""
		n.Children = []*Node{}
	}
	n.Children = append(n.Children, children...)
}

// SetFirst replaces the first child node with the same key as child.
//
// If no child node has the same key, child is appended. Duplicate keys after the
// first match are preserved. It returns true when an existing child was replaced.
func (n *Node) SetFirst(child *Node) bool {
	if n == nil || child == nil {
		return false
	}
	if n.Children == nil {
		n.Value = ""
		n.Children = []*Node{}
	}
	for i, existing := range n.Children {
		if existing != nil && existing.Key == child.Key {
			n.Children[i] = child
			return true
		}
	}
	n.Children = append(n.Children, child)
	return false
}

// RemoveFirst removes and returns the first child node with key.
func (n *Node) RemoveFirst(key string) *Node {
	if n == nil || n.Children == nil {
		return nil
	}
	removed, children := removeFirstNode(n.Children, key)
	n.Children = children
	return removed
}

// RemoveAll removes and returns all child nodes with key.
func (n *Node) RemoveAll(key string) []*Node {
	if n == nil || n.Children == nil {
		return nil
	}
	removed, children := removeAllNodes(n.Children, key)
	n.Children = children
	return removed
}

func cloneNodes(nodes []*Node) []*Node {
	if nodes == nil {
		return nil
	}
	clone := make([]*Node, len(nodes))
	for i, node := range nodes {
		clone[i] = node.Clone()
	}
	return clone
}

func removeFirstNode(nodes []*Node, key string) (*Node, []*Node) {
	for i, node := range nodes {
		if node != nil && node.Key == key {
			out := append(nodes[:i:i], nodes[i+1:]...)
			return node, out
		}
	}
	return nil, nodes
}

func removeAllNodes(nodes []*Node, key string) ([]*Node, []*Node) {
	var removed []*Node
	kept := nodes[:0]
	for _, node := range nodes {
		if node != nil && node.Key == key {
			removed = append(removed, node)
			continue
		}
		kept = append(kept, node)
	}
	if len(removed) == 0 {
		return nil, nodes
	}
	return removed, kept
}
