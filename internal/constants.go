package internal

const (
	TUnknown   = "unknown"
	TInvalid   = "invalid"
	TEof       = "eof"       // \0
	TNewline   = "newline"   // \n
	TKeyword   = "keyword"   // fn, if, else
	TVarIdent  = "varident"  // foo, bar, baz
	TTypeIdent = "typeident" // Foo, Bar, Baz
	TInteger   = "integer"   // 123, 1_000_000
	TFloat     = "float"     // 123.456, 1_000_000.0
	THex       = "hex"       // 0x123
	TOctal     = "octal"     // 0o123
	TBinary    = "binary"    // 0b101
	TString    = "string"    // "foo", `bar`
	TBool      = "bool"      // true, false
	TComment   = "comment"   // --
	TLparen    = "lparen"    // (
	TRparen    = "rparen"    // )
	TLbrace    = "lbrace"    // {
	TRbrace    = "rbrace"    // }
	TLbracket  = "lbracket"  // [
	TRbracket  = "rbracket"  // ]
)

const (
	KImport = "import"
	KFn     = "fn"
)

var Keywords = []string{
	KImport,
	KFn,
}

var TripleCharTokens = map[string]string{
	// "<=>": t_spaceship,
}

var DoubleCharTokens = map[string]string{
	// "==": t_eq,
	// "!=": t_ne,
	// "<=": t_le,
	// ">=": t_ge,
	// "&&": t_and,
	// "||": t_or,
	// "..": t_spread,
	// "->": t_arrow,
	// "+=": t_assigncomp,
	// "-=": t_assigncomp,
	// "*=": t_assigncomp,
	// "/=": t_assigncomp,
	// "%=": t_assigncomp,
}

var SingleCharTokens = map[string]string{
	// "+": t_plus,
	// "-": t_minus,
	// "*": t_star,
	// "/": t_slash,
	// "%": t_percent,
	// "<": t_lt,
	// ">": t_gt,
	// "!": t_bang,
	// "?": t_question,
	// ".": t_dot,
	// ",": t_comma,
	// ":": t_colon,
	// ";": t_semicolon,
	// "=": t_assign,
	// "|": t_pipe,
	"(": TLparen,
	")": TRparen,
	"{": TLbrace,
	"}": TRbrace,
	"[": TLbracket,
	"]": TRbracket,
	// "~": t_tilde,
}
