package token

import "fmt"

type Span struct {
	Filename   string
	FromLine   int
	FromColumn int
	ToLine     int
	ToColumn   int
}

type TokenKind uint64

const (
	TUnknown   TokenKind = iota
	TEof                 // \0
	TNewline             // \n
	TComment             // -- comment
	TSemicolon           // ;

	TConst     // const
	TVarIdent  // variable identifier
	TTypeIdent // type identifier

	// Groupings
	TLeftBrace  // {
	TRightBrace // }

	// Primitive Constants
	TInt    // 0, 1, 2
	THex    // 0x0, 0x1, 0x2
	TOctal  // 0o0, 0o1, 0o2
	TBinary // 0b0, 0b1, 0b2
	TFloat  // 1.0, 1e10, 1.0e10
	TString // 'string'
	TTrue   // true
	TFalse  // false

	// Operators
	TPlus         // +
	TMinus        // -
	TStar         // *
	TSlash        // /
	TPercent      // %
	TGreater      // >
	TGreaterEqual // >=
	TLess         // <
	TLessEqual    // <=
	TSpaceShip    // <=>
	TEqual        // ==
	TNotEqual     // !=
	TAnd          // and
	TOr           // or
	TXor          // xor
	TBang         // !

	// Assignments
	TAssign // =
)

type Token struct {
	Kind  TokenKind
	Loc   Span
	Value string
}

func (t Token) Display() string {
	return KindToLiteral(t.Kind)
}

func (t Token) Is(kind ...TokenKind) bool {
	for _, k := range kind {
		if t.Kind == k {
			return true
		}
	}
	return false
}

var literal2kind = map[string]TokenKind{
	";":     TSemicolon,
	"const": TConst,
	"true":  TTrue,
	"false": TFalse,
	"{":     TLeftBrace,
	"}":     TRightBrace,
	"+":     TPlus,
	"-":     TMinus,
	"*":     TStar,
	"/":     TSlash,
	"%":     TPercent,
	">":     TGreater,
	">=":    TGreaterEqual,
	"<":     TLess,
	"<=":    TLessEqual,
	"<=>":   TSpaceShip,
	"==":    TEqual,
	"!=":    TNotEqual,
	"and":   TAnd,
	"or":    TOr,
	"xor":   TXor,
	"!":     TBang,
	"=":     TAssign,
}

var kind2literal = map[TokenKind]string{
	TUnknown:      "unknown",
	TEof:          "eof",
	TNewline:      "\\n",
	TComment:      "--",
	TSemicolon:    ";",
	TConst:        "const",
	TVarIdent:     "value identifier",
	TTypeIdent:    "type identifier",
	TLeftBrace:    "{",
	TRightBrace:   "}",
	TInt:          "int",
	THex:          "hex",
	TOctal:        "oct",
	TBinary:       "bin",
	TFloat:        "float",
	TString:       "string",
	TTrue:         "true",
	TFalse:        "false",
	TPlus:         "+",
	TMinus:        "-",
	TStar:         "*",
	TSlash:        "/",
	TPercent:      "%",
	TGreater:      ">",
	TGreaterEqual: ">=",
	TLess:         "<",
	TLessEqual:    "<=",
	TSpaceShip:    "<=>",
	TEqual:        "==",
	TNotEqual:     "!=",
	TAnd:          "and",
	TOr:           "or",
	TXor:          "xor",
	TBang:         "!",
	TAssign:       "=",
}

func LiteralToKind(lit string) TokenKind {
	if kind, ok := literal2kind[lit]; ok {
		return kind
	}
	return TUnknown
}

func KindToLiteral(kind TokenKind) string {
	if lit, ok := kind2literal[kind]; ok {
		return lit
	}
	return fmt.Sprintf("not registered '%d'", kind)
}