package yamp

type opType int

// Common operator types
const (
	Addition       opType = iota + 1 // Addition is the binary addition operator.
	Subtraction                      // Subtraction is the binary subtraction operator.
	Multiplication                   // Multiplication is the multiplication operator.
	Division                         // Division is the division operator.
	Power                            // Power is the power/exponent operator.
	Plus                             // Plus is the unary plus sign.
	Minus                            // Minus is the unary minus sign.
	Factorial                        // Factorial is the factorial operator.
)

type assoc int

// Two types of operator associativity
const (
	LeftAssoc  assoc = iota + 1 // LeftAssoc represents the left associativity of an operator.
	RightAssoc                  // RightAssoc represents the right associativity of an operator.
)

type arity int

// Two types of operator arity
const (
	Binary arity = iota + 1 // Binary represents a binary operator.
	Unary                   // Unary represents a unary operator.
)

// Operator represents an arithmetic operator
type Operator interface {
	Token
	Type() opType         // Type returns the operator's type.
	Precedence() int      // Precedence returns the operator's precedence.
	Associativity() assoc // Associativity returns the operator's associativity.
	Arity() arity         // Arity returns the operator's arity.
}

var _ Operator = (*operator)(nil)

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
	return "<?>"
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
	case Plus, Minus:
		return 4
	case Power:
		return 5
	case Factorial:
		return 6
	}
	return 0
}

// Associativity implements the Operator interface.
func (o operator) Associativity() assoc {
	switch o.opType {
	case Addition, Subtraction, Multiplication, Division, Factorial:
		return LeftAssoc
	case Plus, Minus, Power:
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

// IsUnaryOp is a utility function that checks whether an operator
// is unary.
func IsUnaryOp(op Operator) bool {
	return op.Arity() == Unary
}

// IsBinaryOp is a utility function that checks whether an operator
// is binary.
func IsBinaryOp(op Operator) bool {
	return op.Arity() == Binary
}

// IsLeftAssocOp is a utility function that checks whether an operator
// is left associative.
func IsLeftAssocOp(op Operator) bool {
	return op.Associativity() == LeftAssoc
}

// IsRightAssocOp is a utility function that checks whether an operator
// is right associative.
func IsRightAssocOp(op Operator) bool {
	return op.Associativity() == RightAssoc
}
