package lang

type TokenScanner struct {
	*ErrorData
	input  []*Token
	cursor int
	eof    *Token
	queue  []*Token
}

func NewTokenScanner(input []*Token) *TokenScanner {
	return &TokenScanner{
		ErrorData: NewErrorData(),
		input:     input,
		cursor:    0,
		eof:       nil,
		queue:     []*Token{},
	}
}

func (s *TokenScanner) IsEof() bool {
	return s.eof != nil
}

func (s *TokenScanner) EatToken() *Token {
	t := s.PeekToken()
	if len(s.queue) > 0 {
		s.queue = s.queue[1:]
	}
	return t
}

func (s *TokenScanner) EatTokens(n int) []*Token {
	tokens := make([]*Token, 0)
	for i := 0; i < n; i++ {
		tokens = append(tokens, s.EatToken())
	}
	return tokens
}

func (s *TokenScanner) PeekToken() *Token {
	return s.PeekTokenAt(0)
}

func (s *TokenScanner) PeekTokenAt(offset int) *Token {
	for len(s.queue) <= offset {
		token, err := s.next()
		if err != nil {
			s.RegisterError(NewError(
				token.Loc,
				ErrIO,
				err.Error(),
			))
		}
		s.queue = append(s.queue, token)
	}
	return s.queue[offset]
}

func (s *TokenScanner) next() (*Token, error) {
	if s.eof != nil {
		return s.eof, nil
	}
	token := s.input[s.cursor]
	s.cursor++
	if s.cursor >= len(s.input) {
		s.eof = token
	}
	return token, nil
}
