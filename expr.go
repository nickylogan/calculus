package yamp

// Expression represents a mathematical expression
type Expression interface {
	// Evaluate evaluates the expression into a result.
	Evaluate() (res float64, err error)
	String() string
}

var _ Expression = (*expression)(nil)

type expression struct {
	expr   string
	tokens []Token
}

// NewExpression creates a new expression based on the given expression string.
func NewExpression(expr string) Expression {
	return &expression{expr: expr}
}

// Evaluate implements the Expression interface.
func (e *expression) Evaluate() (res float64, err error) {
	t := NewTokenizer()

	tokens, err := t.Tokenize(e.expr)
	if err != nil {
		return 0, err
	}

	rpn := e.toRPN(tokens)
	return e.eval(rpn)
}

type tokenStack struct {
	ts []Token
}

func (s tokenStack) top() Token {
	if len(s.ts) == 0 {
		return nil
	}

	return s.ts[len(s.ts)-1]
}

func (s tokenStack) pop() Token {
	if len(s.ts) == 0 {
		return nil
	}

	t := s.top()
	s.ts = s.ts[:len(s.ts)-1]

	return t
}

func (s tokenStack) push(t Token) {
	s.ts = append(s.ts, t)
}

func (s tokenStack) len() int {
	return len(s.ts)
}

func (e *expression) toRPN(tokens []Token) (rpn []Token) {
	return
}

func (e *expression) eval(rpn []Token) (res float64, err error) {
	return
}

func (e *expression) String() string {
	return e.expr
}
