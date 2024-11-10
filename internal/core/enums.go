package core

// Used to identify the kind of the expression (type or value), useful
// to describe if the construct should be used to evaluate an Type or
// a value
type ExpressionKind string

var (
	InvalidExpression ExpressionKind = "invalid"
	TypeExpression    ExpressionKind = "type"
	ValueExpression   ExpressionKind = "value"
)
