package calculus

type opType int

// Common operator types
const (
	Addition       opType = iota // Addition is the binary addition operator.
	Subtraction                  // Subtraction is the binary subtraction operator.
	Multiplication               // Multiplication is the multiplication operator.
	Division                     // Division is the division operator.
	Power                        // Power is the power/exponent operator.
	Plus                         // Plus is the unary plus sign.
	Minus                        // Minus is the unary minus sign.
	Factorial                    // Factorial is the factorial operator.
)

type assoc int

// Two types of operator associativity
const (
	LeftAssoc  assoc = iota // LeftAssoc represents the left associativity of an operator.
	RightAssoc              // RightAssoc represents the right associativity of an operator.
)

type arity int

// Two types of operator arity
const (
	Binary arity = iota // Binary represents a binary operator.
	Unary               // Unary represents a unary operator.
)

// Operator represents an arithmetic operator
type Operator interface {
	Token
	Type() opType         // Type returns the operator's type.
	Precedence() int      // Precedence returns the operator's precedence.
	Associativity() assoc // Associativity returns the operator's associativity.
	Arity() arity         // Arity returns the operator's arity.
}

type operator struct {
	opType opType
}

// NewOperator creates a new Operator.
func NewOperator(op opType) Operator {
	return operator{opType: op}
}

func (o operator) String() string {
	switch o.opType {
	case Addition, Plus:
		return "+"
	case Subtraction, Minus:
		return "-"
	case Multiplication:
		return "*"
	case Division:
		return "/"
	case Power:
		return "^"
	case Factorial:
		return "!"
	}
	return string(o.opType)
}

// Type implements the Operator interface.
func (o operator) Type() opType {
	return o.opType
}

// Precedence implements the Operator interface.
func (o operator) Precedence() int {
	switch o.opType {
	case Addition, Subtraction:
		return 2
	case Multiplication, Division:
		return 3
	case Power:
		return 4
	case Plus, Minus:
		return 5
	case Factorial:
		return 6
	}
	return 0
}

// Associativity implements the Operator interface.
func (o operator) Associativity() assoc {
	switch o.opType {
	case Addition, Subtraction, Multiplication, Division:
		return LeftAssoc
	case Power, Factorial:
		return RightAssoc
	}
	return LeftAssoc
}

// Arity implements the Operator interface.
func (o operator) Arity() arity {
	switch o.opType {
	case Addition, Subtraction, Multiplication, Division, Power:
		return Binary
	case Plus, Minus, Factorial:
		return Unary
	}
	return Binary
}
