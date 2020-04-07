package calculus

import "strconv"

// Number represents a number
type Number interface {
	Token
	// Value returns the numerical value of the Number
	Value() (res float64, err error)
}

type number struct {
	symbol string
}

func (n number) String() string {
	return n.symbol
}

// Value implements the Number interface
func (n number) Value() (res float64, err error) {
	return strconv.ParseFloat(n.symbol, 64)
}

// NewNumber creates a new number
func NewNumber(num string) Number {
	return number{
		symbol: num,
	}
}
