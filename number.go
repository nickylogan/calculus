package yamp

import (
	"fmt"
	"strconv"
	"strings"
)

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
	res, err = strconv.ParseFloat(n.symbol, 64)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), strconv.ErrSyntax.Error()):
			err = fmt.Errorf(errNaN, n.symbol)
		case strings.Contains(err.Error(), strconv.ErrRange.Error()):
			err = nil
		}
	}
	return
}

// NewNumber creates a new number
func NewNumber(num string) Number {
	return number{
		symbol: num,
	}
}
