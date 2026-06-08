package vdf

// Document is a parsed VDF / KeyValues document.
//
// The top level may contain multiple root nodes. Node order and duplicate keys
// are preserved.
type Document struct {
	Nodes []*Node
}

// NewDocument creates a document from the provided root nodes.
func NewDocument(nodes ...*Node) *Document {
	return &Document{Nodes: nodes}
}

// First returns the first root node with key.
func (d *Document) First(key string) *Node {
	if d == nil {
		return nil
	}
	return firstNode(d.Nodes, key)
}

// All returns all root nodes with key.
func (d *Document) All(key string) []*Node {
	if d == nil {
		return nil
	}
	return allNodes(d.Nodes, key)
}

// Path returns the first node at the provided key path.
func (d *Document) Path(keys ...string) *Node {
	if d == nil || len(keys) == 0 {
		return nil
	}
	node := d.First(keys[0])
	if node == nil {
		return nil
	}
	return node.Path(keys[1:]...)
}
