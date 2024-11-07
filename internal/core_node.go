package internal

import "github.com/renatopp/golden/lang"

type AstData interface {
	Kind() string // type or value
	String() string
}

type RtType interface {
	Name() string
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

func (n *Node) ReplaceBy(node *Node) {
	n.Token = node.Token
	n.Data = node.Data
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
