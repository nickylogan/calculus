package yamp

// Expression represents a mathematical expression
type Expression interface {
	// Evaluate evaluates the expression into a result.
	Evaluate() (res float64, err error)
	String() string
}

type expression struct {
	expr string
}

// NewExpression creates a new expression based on the given expression string.
func NewExpression(expr string) Expression {
	return &expression{expr: expr}
}

// Evaluate implements the Expression interface.
func (e *expression) Evaluate() (res float64, err error) {
	// TODO: implement method
	return
}

func (e *expression) String() string {
	return e.expr
}
