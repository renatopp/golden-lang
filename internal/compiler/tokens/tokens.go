package tokens

const (
	// General
	TUnknown   = "unknown"
	TInvalid   = "invalid"
	TEof       = "eof"        // \0
	TNewline   = "newline"    // \n
	TVarIdent  = "var_ident"  // foo, bar, baz
	TTypeIdent = "type_ident" // Foo, Bar, Baz

	// Primitives
	TInteger = "integer" // 123, 1_000_000
	TFloat   = "float"   // 123.456, 1_000_000.0
	THex     = "hex"     // 0x123
	TOctal   = "octal"   // 0o123
	TBinary  = "binary"  // 0b101
	TString  = "string"  // "foo", `bar`
	TBool    = "bool"    // true, false

	// Keywords
	TImport = "import" // import
	TData   = "data"   // data
	TFn     = "fn"     // fn
	TFN     = "FN"     // FN
	TLet    = "let"    // let
	TAs     = "as"     // as

	// Punctuation
	TComment   = "comment"   // --
	TLparen    = "lparen"    // (
	TRparen    = "rparen"    // )
	TLbrace    = "lbrace"    // {
	TRbrace    = "rbrace"    // }
	TLbracket  = "lbracket"  // [
	TRbracket  = "rbracket"  // ]
	TSemicolon = "semicolon" // ;
	TComma     = "comma"     // ,
	TPipe      = "pipe"      // |
	TDot       = "dot"       // .
	TPlus      = "plus"      // +
	TMinus     = "minus"     // -
	TStar      = "star"      // *
	TSlash     = "slash"     // /
	TPercent   = "percent"   // %
	TLt        = "lt"        // <
	TLte       = "lte"       // <=
	TGt        = "gt"        // >
	TGte       = "gte"       // >=
	TSpaceship = "spaceship" // <=>
	TAssign    = "assign"    // =
	TEqual     = "equal"     // ==
	TNequal    = "not_equal" // !=
	TBang      = "bang"      // !
	TAnd       = "and"       // and
	TOr        = "or"        // or
	TXor       = "xor"       // xor
)

func LiteralToToken(lit string) string {
	switch lit {
	case "import":
		return TImport
	case "data":
		return TData
	case "fn":
		return TFn
	case "FN":
		return TFN
	case "let":
		return TLet
	case "as":
		return TAs
	case "and":
		return TAnd
	case "or":
		return TOr
	case "xor":
		return TXor
	case "true", "false":
		return TBool
	case "<=>":
		return TSpaceship
	case "==":
		return TEqual
	case "!=":
		return TNequal
	case "<=":
		return TLte
	case ">=":
		return TGte
	case "+":
		return TPlus
	case "-":
		return TMinus
	case "*":
		return TStar
	case "/":
		return TSlash
	case "%":
		return TPercent
	case "<":
		return TLt
	case ">":
		return TGt
	case "!":
		return TBang
	case ".":
		return TDot
	case ",":
		return TComma
	case "=":
		return TAssign
	case "|":
		return TPipe
	case ";":
		return TSemicolon
	case "(":
		return TLparen
	case ")":
		return TRparen
	case "{":
		return TLbrace
	case "}":
		return TRbrace
	case "[":
		return TLbracket
	case "]":
		return TRbracket
	}
	return TUnknown
}
