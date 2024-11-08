package internal

import "github.com/renatopp/golden/lang"

type AstData interface {
	Kind() string // type or value
	Label() string
	String() string
	Children() []*Node
}

type RtType interface {
	Name() string
	Accepts(other RtType) bool
	Default() AstData
}

type RtTypeApplicable interface {
	Apply(args []RtType) (RtType, error)
}

type Node struct {
	Token *lang.Token
	Data  AstData
	Type  RtType
}

func NewNode(token *lang.Token, data AstData) *Node {
	return &Node{Token: token, Data: data}
}

func NewEmptyNode() *Node {
	return &Node{}
}

func (n *Node) Copy() *Node {
	return &Node{
		Token: n.Token,
		Data:  n.Data,
		Type:  n.Type,
	}
}

func (n *Node) WithToken(token *lang.Token) *Node {
	n.Token = token
	return n
}

func (n *Node) WithData(data AstData) *Node {
	n.Data = data
	return n
}

func (n *Node) WithType(tp RtType) *Node {
	n.Type = tp
	return n
}

func (n *Node) Traverse(fn func(*Node, int)) {
	n.traverse(fn, 0)
}

func (n *Node) traverse(fn func(*Node, int), depth int) {
	fn(n, depth)
	for _, child := range n.Data.Children() {
		child.traverse(fn, depth+1)
	}
}

func (n *Node) Label() string {
	return f("[%s]", n.Data.Label())
}

func (n *Node) String() string {
	r := ""

	if n.Type != nil {
		tp := n.Type.Name()
		r += f("\n%s::", tp)
	}

	if n.Data != nil {
		r += n.Data.String()
	} else {
		r += "internal"
	}
	return r
}
