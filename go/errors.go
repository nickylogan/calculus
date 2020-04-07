package calculus

const (
	errUnknownSymbol       = "unknown symbol '%s' at index %d"
	errMultipleDecimal     = "cannot allow multiple decimal points in a single number"
	errLoneDecimal         = "something must be on either side of a '.' at index %d"
	errUnmatchedRightParen = "the ')' at index %d is missing a matching '('"
	errUnmatchedLeftParen  = "the '(' at index %d is missing a matching ')'"
	errEmptyParen          = "cannot allow an empty parentheses on index %d"
	errNoRightOperand      = "operator '%s' at index %d expects a right operand"
	errNoLeftOperand       = "operator '%s' at index %d requires a left operand"
)

// SyntaxError stores a syntax error.
type SyntaxError struct {
	Message  string
	Token    string
	Position int
}

func (s SyntaxError) Error() string {
	return s.Message
}
