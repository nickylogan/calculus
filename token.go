package calculus

import "unicode"

// Token represents a single expression token
type Token interface {
	String() string
}

// IsDigit checks if a given rune is a digit
func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// IsOperator checks if a given rune is an operator
func IsOperator(r rune) bool {
	return r == '+' ||
		r == '-' ||
		r == '*' ||
		r == '/' ||
		r == '^' ||
		r == '!'
}

// IsLeftParen checks if a given rune is a left parenthesis.
func IsLeftParen(r rune) bool {
	return r == '('
}

// IsRightParen checks if a given rune is a right parenthesis.
func IsRightParen(r rune) bool {
	return r == ')'
}

// IsDecimalPoint checks if a given rune is a decimal point.
func IsDecimalPoint(r rune) bool {
	return r == '.'
}

// IsWhitespace checks if a given rune is a whitespace.
func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
