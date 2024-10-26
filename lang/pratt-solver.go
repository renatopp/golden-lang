package lang

type PrattSolver[T comparable] struct {
	precedence func(*Token) int
	prefixFns  map[string]func() T
	infixFns   map[string]func(T) T
	postfixFns map[string]func(T) T
}

func NewPrattSolver[T comparable]() *PrattSolver[T] {
	return &PrattSolver[T]{
		precedence: func(t *Token) int { return 0 },
		prefixFns:  make(map[string]func() T),
		infixFns:   make(map[string]func(T) T),
		postfixFns: make(map[string]func(T) T),
	}
}

func (p *PrattSolver[T]) RegisterPrefixFn(kind string, fn func() T) *PrattSolver[T] {
	p.prefixFns[kind] = fn
	return p
}

func (p *PrattSolver[T]) RegisterInfixFn(kind string, fn func(T) T) *PrattSolver[T] {
	p.infixFns[kind] = fn
	return p
}

func (p *PrattSolver[T]) RegisterPostfixFn(kind string, fn func(T) T) *PrattSolver[T] {
	p.postfixFns[kind] = fn
	return p
}

func (p *PrattSolver[T]) SetPrecedence(fn func(*Token) int) *PrattSolver[T] {
	p.precedence = fn
	return p
}

func (p *PrattSolver[T]) SolveExpression(s *TokenScanner, precedence int) T {
	var zero T
	prefix := p.prefixFns[s.PeekToken().Kind]
	if prefix == nil {
		return zero
	}
	left := prefix()
	if left == zero {
		return zero
	}

	cur := s.PeekToken()
	for {
		starting := cur

		if precedence < p.precedence(cur) {
			infix := p.infixFns[cur.Kind]
			if infix != nil {
				left = infix(left)
				cur = s.PeekToken()
			}
		}

		for {
			postfix := p.postfixFns[cur.Kind]
			if postfix == nil {
				break
			}
			newLeft := postfix(left)
			if newLeft == zero {
				break
			}
			left = newLeft
			cur = s.PeekToken()
		}

		// Didn't find any infix or postfix function
		if starting == cur {
			break
		}
	}

	return left
}
