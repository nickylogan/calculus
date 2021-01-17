package yamp

import "unicode"

// IsDigit checks if a given rune is a digit.
func IsDigit(r rune) bool {
	return defaultTokenRegistry.IsDigit(r)
}

// IsOperator checks if a given rune is an operator.
func IsOperator(r rune) bool {
	return defaultTokenRegistry.IsOperator(r)
}

// IsLeftBracket checks if a given rune is a left parenthesis.
func IsLeftBracket(r rune) bool {
	return defaultTokenRegistry.IsLeftBracket(r)
}

// IsRightBracket checks if a given rune is a right parenthesis.
func IsRightBracket(r rune) bool {
	return defaultTokenRegistry.IsRightBracket(r)
}

// IsDecimalPoint checks if a given rune is a decimal point.
func IsDecimalPoint(r rune) bool {
	return defaultTokenRegistry.IsDecimalPoint(r)
}

// IsWhitespace checks if a given rune is a whitespace.
func IsWhitespace(r rune) bool {
	return defaultTokenRegistry.IsWhitespace(r)
}

var defaultTokenRegistry = TokenRegistry{
	operators: defaultOperators,
}

// TokenRegistry maps runes into their respective tokens.
type TokenRegistry struct {
	operators OperatorRegistry
}

// IsDigit checks if a given rune is a digit.
func (m TokenRegistry) IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// IsOperator checks if a given rune is an operator.
func (m TokenRegistry) IsOperator(r rune) bool {
	_, ok := m.operators.GetOperator(r, nil)
	return ok
}

// IsLeftBracket checks if a given rune is a left parenthesis.
func (m TokenRegistry) IsLeftBracket(r rune) bool {
	// TODO: add support for other bracket notations.
	return r == '('
}

// IsRightBracket checks if a given rune is a right parenthesis.
func (m TokenRegistry) IsRightBracket(r rune) bool {
	// TODO: add support for other bracket notations.
	return r == ')'
}

// IsDecimalPoint checks if a given rune is a decimal point.
func (m TokenRegistry) IsDecimalPoint(r rune) bool {
	// TODO: add support for other decimal point notations.
	return r == '.'
}

// IsWhitespace checks if a given rune is a whitespace.
func (m TokenRegistry) IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

var defaultOperators = OperatorRegistry{
	'+': {NewOperator(Plus), NewOperator(Addition)},
	'-': {NewOperator(Minus), NewOperator(Subtraction)},
	'*': {NewOperator(Multiplication)},
	'/': {NewOperator(Division)},
	'^': {NewOperator(Power)},
	'!': {NewOperator(Factorial)},
}

type (
	// OperatorRegistry contains a registry of operator runes.
	// Use make(OperatorRegistry) to create a new OperatorRegistry.
	OperatorRegistry map[rune][]Operator
)

// Register registers a new rune with the given operator token.
func (reg OperatorRegistry) Register(r rune, op Operator) {
	reg[r] = append(reg[r], op)
}

// GetOperator gets the operator of a rune that matches the given filter(s). Filters are combined using the AND clause.
// If no filters are given, GetOperator will return the first operator of that rune it encounters in the
// registry. If none is found, ok will be false.
func (reg OperatorRegistry) GetOperator(r rune, filter ...func(Operator) bool) (res Operator, ok bool) {
Loop:
	for _, op := range reg[r] {
		if len(filter) == 0 {
			return op, true
		}

		for _, f := range filter {
			if f == nil {
				continue
			}
			if !f(op) {
				continue Loop
			}
		}

		return op, true
	}

	return nil, false
}
