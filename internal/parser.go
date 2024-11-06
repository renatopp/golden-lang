package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token) (*Node, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		Parser: lang.NewParser(scanner),
		Pratt:  lang.NewPrattSolver[*Node](),
	}

	parser.Pratt.SetPrecedenceFn(parser.precedence)
	parser.Pratt.RegisterPrefixFn(TInteger, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(THex, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TOctal, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TBinary, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TFloat, parser.parseFloat)
	parser.Pratt.RegisterPrefixFn(TBool, parser.parseBool)
	parser.Pratt.RegisterPrefixFn(TString, parser.parseString)
	parser.Pratt.RegisterPrefixFn(TVarIdent, parser.parseVarIdent)
	parser.Pratt.RegisterPrefixFn(TTypeIdent, parser.parseTypeCall)
	parser.Pratt.RegisterPrefixFn(TLbrace, parser.parseBlock)
	parser.Pratt.RegisterPrefixFn(TPlus, parser.parseUnaryOperator)
	parser.Pratt.RegisterPrefixFn(TMinus, parser.parseUnaryOperator)
	parser.Pratt.RegisterPrefixFn(TBang, parser.parseUnaryOperator)
	parser.Pratt.RegisterPrefixFn(TLparen, parser.parseAnonymousType)

	parser.Pratt.RegisterInfixFn(TPlus, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TMinus, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TStar, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TSlash, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TPercent, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TSpaceship, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TEqual, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TNequal, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TAnd, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TOr, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TXor, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TLt, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TLte, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TGt, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TGte, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TLparen, parser.parseFunctionCall)
	parser.Pratt.RegisterInfixFn(TDot, parser.parseAccess)

	node := parser.Parse()
	if parser.Scanner.HasErrors() || parser.HasErrors() {
		return nil, lang.NewErrorList(append(parser.Errors(), parser.Scanner.Errors()...))
	}

	return node, nil
}

type parser struct {
	*lang.Parser
	Pratt *lang.PrattSolver[*Node]
}

func (p *parser) Parse() (out *Node) {
	defer func() {
		r := recover()
		if r == nil {
			return
		} else if err, ok := r.(lang.Error); ok {
			p.RegisterError(err)
		} else {
			p.RegisterError(lang.NewError(lang.Loc{}, "unknown error", fmt.Sprintf("%v", r)))
		}
	}()

	return p.parseModule()
}

// Overriding methods from lang.Parser

func (p *parser) ExpectToken(kinds ...string) {
	if !p.Parser.ExpectToken(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectLiteral(literals ...string) {
	if !p.Parser.ExpectLiteral(literals...) {
		panic(nil)
	}
}
func (p *parser) Expect(kind string, literals ...string) {
	if !p.Parser.Expect(kind, literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipToken1(kinds ...string) {
	if !p.Parser.ExpectSkipToken1(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipTokenAll(kinds ...string) {
	if !p.Parser.ExpectSkipTokenAll(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipLiteral1(literals ...string) {
	if !p.Parser.ExpectSkipLiteral1(literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipLiteralAll(literals ...string) {
	if !p.Parser.ExpectSkipLiteralAll(literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkip1(kind string, literals ...string) {
	if !p.Parser.ExpectSkip1(kind, literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipAll(kind string, literals ...string) {
	if !p.Parser.ExpectSkipAll(kind, literals...) {
		panic(nil)
	}
}

// Custom methods

func (p *parser) precedence(t *lang.Token) int {
	switch {
	case t.IsKind(TAssign):
		return 10
	case t.IsKind(TPipe):
		return 20
	case t.IsKind(TOr):
		return 40
	case t.IsKind(TXor):
		return 45
	case t.IsKind(TAnd):
		return 50
	case t.IsKind(TEqual, TNequal):
		return 70
	case t.IsKind(TLt, TGt, TLte, TGte):
		return 80
	case t.IsKind(TPlus, TMinus):
		return 90
	case t.IsKind(TStar, TSlash):
		return 100
	case t.IsKind(TSpaceship):
		return 110
	case t.IsKind(TPercent):
		return 120
	case t.IsKind(TLparen):
		return 130
	case t.IsKind(TDot):
		return 140
	}

	return 0
}

// Parse the module root node, i.e., the first depth level of the file.
// The module node contains all the imports and declaration (types, variables and functions).
//
// Example:
//
//			import "foo"
//			data Node = Nil | Node
//	   let x Int = 10
//	   fn main() { }
//
// Returns:
// - Node(AstModule)
func (p *parser) parseModule() *Node {
	first := p.PeekToken()
	imports := []*Node{}
	types := []*Node{}
	functions := []*Node{}
	variables := []*Node{}

	for {
		p.Skip(TSemicolon, TNewline)

		switch {
		case p.IsNext(TKeyword, KImport):
			imp := p.parseImport()
			imports = append(imports, imp)
			continue
		case p.IsNext(TKeyword, KData):
			data := p.parseDataDecl()
			types = append(types, data)
			continue
		case p.IsNext(TKeyword, KFn):
			fn := p.parseFunctionDecl()
			functions = append(functions, fn)
			continue
		case p.IsNext(TKeyword, KLet):
			variable := p.parseVariableDecl()
			variables = append(variables, variable)
			continue
		default:
			p.Expect(TEof)
		}

		break
	}

	return NewNode(first, &AstModule{
		Imports:   imports,
		Types:     types,
		Functions: functions,
		Variables: variables,
	})
}

// Parse the import expression with an optional alias.
//
// Example:
//
//	import "foo"
//	import "foo" as bar
//
// Returns:
// - Node(AstImport)
func (p *parser) parseImport() *Node {
	p.Expect(TKeyword, KImport)
	imp := p.EatToken()

	p.Expect(TString)
	path := p.EatToken()

	alias := ""
	if p.IsNext(TKeyword, KAs) {
		p.EatToken()
		p.Expect(TVarIdent)
		alias = p.EatToken().Literal
	}

	return NewNode(imp, &AstImport{
		Alias: alias,
		Path:  path.Literal,
	})
}

// Parse the data declaration with all possible variations.
//
// Example:
//
//	data NodeConstructor
//	data Node = Constructor
//	data Node = Constructor | Constructor | ...
//
// see parseConstructor for more details on the constructor syntax.
//
// Returns:
// - Node(AstDataDecl)
func (p *parser) parseDataDecl() *Node {
	p.Expect(TKeyword, KData)
	data := p.EatToken()

	p.Expect(TTypeIdent)
	name := p.EatToken()

	constructors := []*Node{}

	if p.IsNext(TAssign) {
		// Parse data declaration with explicit constructors
		// Example:
		//
		// 		data Node = Nil | Node
		// 		data Node = Node
		// 		data Node = Node(Int)
		// 		data Node = Node(x Int)
		//
		p.EatToken()
		constructors = p.parseConstructors()

	} else if p.IsNext(TLparen) {
		// Parse data declaration with implicit, same-name constructors
		// Example:
		//
		// 		data Node()
		// 		data Node(Int)
		// 		data Node(x Int)
		//
		c := p.parseConstructor()
		if c == nil {
			panic(lang.NewError(p.PeekToken().Loc, "expecting constructor", ""))
		}
		c.Data.(*AstConstructor).Name = name.Literal
		constructors = append(constructors, c)

	} else {
		// Parse data declaration with implicit, same-name, unit constructor
		// Example:
		//
		// 		data Node
		//
		constructors = append(constructors, NewNode(name, &AstConstructor{
			Name:   name.Literal,
			Shape:  "unit",
			Fields: []*Node{},
		}))
	}

	return NewNode(data, &AstDataDecl{
		Name:         name.Literal,
		Constructors: constructors,
	})
}

// Parse the function declaration with parameters, return type and body.
//
// Example:
//
//	fn main() { }
//	fn add(a Int, b Int) Int { a + b }
//
// Returns:
// - Node(AstFunctionDecl)
func (p *parser) parseFunctionDecl() *Node {
	p.Expect(TKeyword, KFn)
	fn := p.EatToken()

	name := ""
	if p.IsNext(TVarIdent) {
		name = p.EatToken().Literal
	}

	p.Expect(TLparen)
	p.EatToken()

	parameters := p.parseParameters()

	p.Expect(TRparen)
	p.EatToken()

	tp := p.parseTypeRef()

	body := p.parseBlock()

	return NewNode(fn, &AstFunctionDecl{
		Name:       name,
		Parameters: parameters,
		ReturnType: tp,
		Body:       body,
	})
}

// Parse the variable declaration with all variations.
//
// Example:
//
//	let x = 10
//	let x Int
//	let x Int = 10
//
// Returns:
// - Node(AstVariableDecl)
func (p *parser) parseVariableDecl() *Node {
	p.Expect(TKeyword, KLet)
	let := p.EatToken()

	p.ExpectToken(TVarIdent)
	name := p.EatToken()

	tp := p.parseTypeRef()
	if tp == nil {
		p.ExpectToken(TAssign)
	}

	var value *Node
	if p.IsNextToken(TAssign) {
		p.EatToken()
		value = p.parseExpression()
		if value == nil {
			panic(lang.NewError(p.PeekToken().Loc, "expecting expression", ""))
		}
	}

	return NewNode(let, &AstVariableDecl{
		Name:       name.Literal,
		Type:       tp,
		Expression: value,
	})
}

// Parse the constructor list of a data declaration.
//
// Example:
//
//	Constructor
//	Constructor | Constructor | ...
//
// Returns:
// - []*Node(AstConstructor)
func (p *parser) parseConstructors() []*Node {
	constructors := []*Node{}

	// ignore initial newline
	if p.IsNextToken(TNewline) {
		p.Skip(TNewline)

		// ignore initial pipe if newline is present
		if p.IsNext(TPipe) {
			p.EatToken()
			p.Expect(TTypeIdent)
		}
	}

	// parse constructors
	for {
		if !p.IsNext(TTypeIdent) {
			break
		}

		constructor := p.parseConstructor()
		if constructor == nil {
			break
		}
		constructors = append(constructors, constructor)

		p.Skip(TNewline)
		if !p.IsNext(TPipe) {
			break
		}
		p.EatToken()
		p.Expect(TTypeIdent)
	}

	return constructors
}

// Parse a single constructor of a data declaration. It may or may not have a name.
//
// Example:
//
//	Constructor
//	Constructor(Int)
//	Constructor(x Int)
//	()
//	(Int)
//	(x Int)
//
// Returns:
// - Node(AstConstructor)
func (p *parser) parseConstructor() *Node {
	first := p.PeekToken()

	name := ""
	if p.IsNext(TTypeIdent) {
		name = p.EatToken().Literal
	}

	p.Skip(TNewline)
	if !p.IsNext(TLparen) {
		return NewNode(first, &AstConstructor{
			Name:   name,
			Shape:  "unit",
			Fields: []*Node{},
		})
	}

	p.EatToken() // (
	p.Skip(TNewline)

	shape, fields := p.parseConstructorFields()

	p.Skip(TNewline)
	p.Expect(TRparen)
	p.EatToken()
	return NewNode(first, &AstConstructor{
		Name:   name,
		Shape:  shape,
		Fields: fields,
	})
}

// Parse the fields of a constructor, which can treated as an unit,  a tuple or a record.
// Notice that this functions considers the content INSIDE the parenthesis.
//
// Example:
//
//	<empty>
//	Int, ...
//	x Int, ...
//
// Returns:
// - shape, []*Node(AstField)
func (p *parser) parseConstructorFields() (shape string, fields []*Node) {
	shape = "unit" // or "tuple" or "record"
	fields = []*Node{}

	// discover the shape of the constructor
	p.Skip(TNewline)
	switch {
	case p.IsNext(TVarIdent):
		shape = "record"
		for !p.IsNext(TRparen) {
			name := p.EatToken()
			tp := p.parseTypeRef()
			if tp == nil {
				panic(lang.NewError(p.PeekToken().Loc, "expecting type", ""))
			}
			fields = append(fields, NewNode(name, &AstField{
				Name: name.Literal,
				Type: tp,
			}))

			p.ExpectToken(TComma, TNewline, TRparen)
			p.Skip(TNewline)
			p.SkipN(1, TComma)
			p.Skip(TNewline)
		}

	case p.IsNext(TRparen):
		shape = "unit"

	default:
		shape = "tuple"
		i := 0
		for !p.IsNext(TRparen) {
			tp := p.parseTypeRef()
			if tp == nil {
				panic(lang.NewError(p.PeekToken().Loc, "expecting type", ""))
			}
			name := strconv.Itoa(i)
			fields = append(fields, NewNode(p.PeekToken(), &AstField{
				Name: name,
				Type: tp,
			}))

			p.ExpectToken(TComma, TNewline, TRparen)
			p.Skip(TNewline)
			p.SkipN(1, TComma)
			p.Skip(TNewline)
			i++
		}
	}

	return shape, fields
}

// Parse the parameters of a function declaration.
// This functions considers the content INSIDE the parenthesis.
//
// Example:
//
//	<empty>
//	x Int, ...
//
// Returns:
// - []*Node(AstParameter)
func (p *parser) parseParameters() []*Node {
	parameters := []*Node{}
	p.Skip(TNewline)
	for {
		if !p.IsNext(TVarIdent) {
			break
		}

		name := p.EatToken()
		tp := p.parseTypeRef()
		parameters = append(parameters, NewNode(name, &AstParameter{
			Name: name.Literal,
			Type: tp,
		}))

		p.ExpectToken(TComma, TNewline, TRparen)
		p.Skip(TNewline)
		p.SkipN(1, TComma)
		p.Skip(TNewline)
	}

	return parameters
}

// Parse the parameters of a function type reference.
// This functions considers the content INSIDE the parenthesis.
//
// Example:
//
//	<empty>
//	Int, ...
//
// Returns:
// - []*Node(AstParameter)
func (p *parser) parseParameterTypes() []*Node {
	parameters := []*Node{}
	p.Skip(TNewline)
	for {
		if !p.IsNext(TTypeIdent) {
			break
		}

		tp := p.parseTypeRef()
		parameters = append(parameters, NewNode(tp.Token, &AstParameter{
			Name: "",
			Type: tp,
		}))
		p.ExpectToken(TComma, TNewline, TRparen)
		p.Skip(TNewline)
		p.SkipN(1, TComma)
		p.Skip(TNewline)
	}

	return parameters
}

// Parse the reference for a type, i.e., a type expression such as the return type of a function or
// the variable type in a declaration.
//
// The type reference can be a named type, a function type or a tuple/record type.
//
// IMPORTANT: It may return nil
//
// Example:
//
//	Int
//	Fn(Int) Int
//	()
//	(Int, Int)
//	(x Int, y Int)
//
// Returns:
// - nil
// - Node(AstTypeRef)
// - Node(AstFnTypeRef)
// - Node(AstDataDecl)
func (p *parser) parseTypeRef() *Node {
	switch {
	case p.IsNext(TLparen):
		// tuple or record
		first := p.EatToken()
		p.Skip(TNewline)

		shape, fields := p.parseConstructorFields()

		// unit is equivalent to void
		if shape == "unit" {
			p.Expect(TRparen)
			p.EatToken()
			return nil
		}

		p.Expect(TRparen)
		p.EatToken()
		return NewNode(first, &AstDataDecl{
			Name: "",
			Constructors: []*Node{
				NewNode(first, &AstConstructor{
					Name:   "",
					Shape:  shape,
					Fields: fields,
				}),
			},
		})

	case p.IsNext(TTypeIdent, "Fn"):
		// function type
		first := p.EatToken()
		p.Expect(TLparen)
		p.EatToken()

		parameters := p.parseParameterTypes()

		p.Skip(TNewline)
		p.Expect(TRparen)
		p.EatToken()

		returnType := p.parseTypeRef()
		return NewNode(first, &AstFnTypeRef{
			Parameters: parameters,
			ReturnType: returnType,
		})

	case p.IsNext(TTypeIdent):
		// named type
		name := p.EatToken()
		return NewNode(name, &AstTypeRef{
			Name: name.Literal,
		})
	}

	return nil
}

// Parse a block expression, i.e., a sequence of expressions enclosed by braces.
//
// Example:
//
//		{
//			let x = 10
//			let y = 20
//	 }
//
// Returns:
// - Node(AstBlock)
func (p *parser) parseBlock() *Node {
	p.Expect(TLbrace)
	first := p.EatToken()
	p.Skip(TNewline, TSemicolon)

	expressions := []*Node{}
	for {
		expr := p.parseExpression()
		if expr == nil {
			break
		}

		expressions = append(expressions, expr)
		p.ExpectToken(TNewline, TSemicolon, TRbrace)
		p.Skip(TNewline, TSemicolon)
	}

	p.Expect(TRbrace)
	p.EatToken()

	return NewNode(first, &AstBlock{
		Expressions: expressions,
	})
}

// Parse a single expression.
//
// ATTENTION: this functions may return nil if not expression is found.
//
// Example:
//
//			10
//			10 + 20
//			call()
//	   let x = 10
//			...
//
// Returns:
// - nil
// - Node(<multiple>)
func (p *parser) parseExpression(precedence ...int) *Node {
	switch {
	case p.IsNext(TKeyword, KLet):
		return p.parseVariableDecl()
	case p.IsNext(TKeyword, KFn):
		return p.parseFunctionDecl()
	}

	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.Pratt.SolveExpression(p.Scanner, pr)
}

// Parse an integer literal with support for different bases.
func (p *parser) parseInteger() *Node {
	p.ExpectToken(TInteger, THex, TOctal, TBinary)

	token := p.EatToken()
	base := 10
	switch token.Kind {
	case THex:
		base = 16
	case TOctal:
		base = 8
	case TBinary:
		base = 2
	}

	value, err := strconv.ParseInt(token.Literal, base, 64)
	if err != nil {
		panic(lang.NewError(token.Loc, "invalid integer", token.Literal))
	}

	return NewNode(token, &AstInt{Value: value})
}

// Parse a float literal.
func (p *parser) parseFloat() *Node {
	p.ExpectToken(TFloat)
	token := p.EatToken()
	value, err := strconv.ParseFloat(token.Literal, 64)
	if err != nil {
		panic(lang.NewError(token.Loc, "invalid float", token.Literal))
	}
	return NewNode(token, &AstFloat{Value: value})
}

// Parse a boolean literal.
func (p *parser) parseBool() *Node {
	p.ExpectToken(TBool)
	p.ExpectLiteral("true", "false")
	token := p.EatToken()
	value := token.Literal == "true"
	return NewNode(token, &AstBool{Value: value})
}

// Parse a string literal.
func (p *parser) parseString() *Node {
	p.ExpectToken(TString)
	token := p.EatToken()
	value := strings.ReplaceAll(token.Literal, "\r", "")
	return NewNode(token, &AstString{Value: value})
}

// Parse a variable identifier.
func (p *parser) parseVarIdent() *Node {
	p.ExpectToken(TVarIdent)
	token := p.EatToken()
	return NewNode(token, &AstVarIdent{Name: token.Literal})
}

func (p *parser) parseTypeCall() *Node {
	p.ExpectToken(TTypeIdent)
	first := p.PeekToken()
	tp := p.parseTypeRef()

	shape := "unit"
	args := []*Node{}
	if p.IsNext(TLparen) {
		p.EatToken()
		shape, args = p.parseArguments()
		p.Expect(TRparen)
		p.EatToken()
	}
	return NewNode(first, &AstTypeCall{
		Shape: shape,
		Type:  tp,
		Args:  args,
	})
}

// Parse a unary operator expression.
func (p *parser) parseUnaryOperator() *Node {
	op := p.EatToken()
	right := p.parseExpression(p.precedence(op))
	return NewNode(op, &AstUnary{
		Op:    op.Literal,
		Right: right,
	})
}

// Parse an anonymous type expression.
func (p *parser) parseAnonymousType() *Node {
	p.Expect(TLparen)
	first := p.EatToken()
	shape, args := p.parseArguments()
	p.Expect(TRparen)
	p.EatToken()
	return NewNode(first, &AstTypeCall{
		Shape: shape,
		Type:  nil, // anonymous
		Args:  args,
	})
}

// Parse a list of expressions in a record format or tuple format.
//
// Example:
//
//	   10, 20
//			x=10, y=20
//
// Returns:
// - Node([]*Ast)
func (p *parser) parseArguments() (string, []*Node) {
	args := []*Node{}
	p.Skip(TNewline)

	shape := ""
	switch {
	case p.IsNext(TVarIdent):
		// record
		shape = "record"
		for {
			p.ExpectToken(TVarIdent)
			name := p.EatToken()
			p.Expect(TAssign)
			p.EatToken()
			p.Skip(TNewline)
			value := p.parseExpression()
			if value == nil {
				panic(lang.NewError(p.PeekToken().Loc, "expecting expression", ""))
			}
			args = append(args, NewNode(name, &AstArgument{
				Name:       name.Literal,
				Expression: value,
			}))

			p.ExpectToken(TComma, TNewline, TRparen)
			p.Skip(TNewline)
			p.SkipN(1, TComma)
			p.Skip(TNewline)
			if p.IsNext(TRparen) {
				break
			}
		}

	default:
		// tuple
		shape = "tuple"
		i := 0
		p.Skip(TNewline)
		for {
			value := p.parseExpression()
			if value == nil {
				break
			}
			args = append(args, NewNode(p.PeekToken(), &AstArgument{
				Name:       strconv.Itoa(i),
				Expression: value,
			}))

			p.ExpectToken(TComma, TNewline, TRparen)
			p.Skip(TNewline)
			p.SkipN(1, TComma)
			p.Skip(TNewline)
			if p.IsNext(TRparen) {
				break
			}
			i++
		}
	}

	if len(args) == 0 {
		shape = "unit"
	}

	return shape, args
}

// Parse a binary operator expression.
func (p *parser) parseBinaryOperator(left *Node) *Node {
	op := p.EatToken()
	right := p.parseExpression(p.precedence(op))

	if right == nil {
		panic(lang.NewError(op.Loc, "expecting expression", ""))
	}

	return NewNode(op, &AstBinary{
		Op:    op.Literal,
		Left:  left,
		Right: right,
	})
}

// Parse a function call or type initialization.
func (p *parser) parseFunctionCall(left *Node) *Node {
	p.Expect(TLparen)
	first := p.EatToken()
	_, args := p.parseArguments()
	p.Expect(TRparen)
	p.EatToken()
	return NewNode(first, &AstFnCall{
		Target: left,
		Args:   args,
	})
}

// Parse member access.
//
// Example:
//
//	foo.bar
//	foo.0
//
// Returns:
// - Node(AstAccess)
func (p *parser) parseAccess(left *Node) *Node {
	p.Expect(TDot)
	first := p.EatToken()
	p.ExpectToken(TVarIdent, TInteger)
	token := p.EatToken()
	return NewNode(first, &AstAccess{
		Target: left,
		Member: token.Literal,
	})
}
