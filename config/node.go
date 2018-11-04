package config

type Node struct {
	Path     string
	Value    string
	Children []*Node
}

func newNode() *Node {
	var node = &Node{
		Path:  "/",
		Value: "",
	}

	return node
}

func (n *Node) appendChild(node *Node) {
	if n.Children == nil {
		n.Children = []*Node{}
	}

	n.Children = append(n.Children, node)
}
