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
	TSemicolon = "semicolon" // ;
	TComma     = "comma"     // ,
	TPipe      = "pipe"      // |

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
	TNequal    = "nequal"    // !=
	TBang      = "bang"      // !
	TAnd       = "and"       // and
	TOr        = "or"        // or
	TXor       = "xor"       // xor
)

const (
	KImport = "import"
	KData   = "data"
	KFn     = "fn"
	KLet    = "let"
)

var Keywords = []string{
	KImport,
	KData,
	KFn,
	KLet,
}

var KeywordTokens = map[string]string{
	"and": TAnd,
	"or":  TOr,
	"xor": TXor,
}

var TripleCharTokens = map[string]string{
	"<=>": TSpaceship,
}

var DoubleCharTokens = map[string]string{
	"==": TEqual,
	"!=": TNequal,
	"<=": TLte,
	">=": TGte,
	// "..": TSpread,
	// "->": TArrow,
	// "+=": TAssign,
	// "-=": TAssigncomp,
	// "*=": TAssigncomp,
	// "/=": TAssigncomp,
	// "%=": TAssigncomp,
}

var SingleCharTokens = map[string]string{
	"+": TPlus,
	"-": TMinus,
	"*": TStar,
	"/": TSlash,
	"%": TPercent,
	"<": TLt,
	">": TGt,
	"!": TBang,
	// "?": TQuestion,
	// ".": TDot,
	// ":": TColon,
	",": TComma,
	"=": TAssign,
	"|": TPipe,
	";": TSemicolon,
	"(": TLparen,
	")": TRparen,
	"{": TLbrace,
	"}": TRbrace,
	"[": TLbracket,
	"]": TRbracket,
	// "~": t_tilde,
}
