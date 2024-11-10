package core

type AstData interface {
	ExpressionKind() ExpressionKind // Expression information (type or value)
	Tag() string                    // Used as a short name identifier of the node for debugging
	Signature() string              // Signature of the expression, used for errors and other information for the user
	Children() []*AstNode           // Nested nodes
}
