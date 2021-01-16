package calculus

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Tokenizer is implemented by an expression tokenizer.
type Tokenizer interface {
	// Tokenize is called when an expression needs to be split into understandable tokens.
	Tokenize(expr string) (tokens []Token, err error)
}

const (
	tokenNothing = 1 << iota
	tokenInteger
	tokenLoneDecimalPoint
	tokenDecimal
	tokenLeftParen
	tokenRightParen
	tokenBinaryOp
	tokenLeftUnaryOp
	tokenRightUnaryOp
)

type tokenizer struct {
	expr       string
	sc         *bufio.Scanner
	tokens     []Token
	currState  int
	currSymbol strings.Builder
	currIndex  int
	parenDepth *depthStack
}

// NewTokenizer creates a new Tokenizer
func NewTokenizer() Tokenizer {
	return &tokenizer{}
}

// Tokenize implements the Tokenizer interface
func (t *tokenizer) Tokenize(expr string) (tokens []Token, err error) {
	t.initialize(expr)
	for t.sc.Scan() {
		if err = t.sc.Err(); err != nil {
			return
		}

		r := []rune(t.sc.Text())[0]
		switch {
		case IsDigit(r):
			if err = t.handleDigit(r); err != nil {
				return
			}
		case IsDecimalPoint(r):
			if err = t.handleDecimalPoint(r); err != nil {
				return
			}
		case IsLeftParen(r):
			if err = t.handleLeftParen(r); err != nil {
				return
			}
		case IsRightParen(r):
			if err = t.handleRightParen(r); err != nil {
				return
			}
		case IsOperator(r):
			if err = t.handleOperator(r); err != nil {
				return
			}
		case IsWhitespace(r):
			// ignore whitespace
		default:
			err = SyntaxError{
				Message:  fmt.Sprintf(errUnknownSymbol, string(r), t.currIndex),
				Token:    string(r),
				Position: t.currIndex,
			}
			return
		}
		t.currIndex++
	}

	// commit last token
	t.commitCurrentState()

	tokens = t.tokens
	return
}

func (t *tokenizer) handleDigit(r rune) (err error) {
	// "." => ".5"
	if t.currState == tokenLoneDecimalPoint {
		t.currSymbol.WriteRune(r)
		t.currState = tokenDecimal
	}
	// "1" => "12", "1.2" => "1.23"
	if t.currState&(tokenInteger|tokenDecimal) != 0 {
		t.currSymbol.WriteRune(r)
		return
	}

	t.commitCurrentState()

	if t.currState&(tokenRightParen|tokenRightUnaryOp) != 0 {
		t.appendToken(NewOperator(Multiplication))
	}
	t.currSymbol.WriteRune(r)
	t.currState = tokenInteger
	return
}

func (t *tokenizer) handleDecimalPoint(r rune) (err error) {
	// can't allow multiple decimal points
	if t.currState&(tokenDecimal|tokenLoneDecimalPoint) != 0 {
		return SyntaxError{
			Message:  errMultipleDecimal,
			Token:    string(r),
			Position: t.currIndex,
		}
	}

	// "5" => "5."
	if t.currState == tokenInteger {
		t.currSymbol.WriteRune(r)
		t.currState = tokenDecimal
		return
	}

	t.commitCurrentState()

	// "(5)" => "(5)*.", "5!" => "5!*."
	if t.currState&(tokenRightParen|tokenRightUnaryOp) != 0 {
		t.appendToken(NewOperator(Multiplication))
	}
	t.currSymbol.WriteRune(r)
	t.currState = tokenLoneDecimalPoint
	return
}

func (t *tokenizer) handleLeftParen(r rune) (err error) {
	// can't allow lone decimal point to be followed by a left parenthesis
	if t.currState == tokenLoneDecimalPoint {
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
		}
	}

	t.commitCurrentState()

	// "5" => "5*(", "5.4" => "5.4*(", "(5)" => "(5)*(", "5!" => "5!*("
	if t.currState&(tokenInteger|tokenDecimal|tokenRightParen|tokenRightUnaryOp) != 0 {
		t.appendToken(NewOperator(Multiplication))
	}

	t.currSymbol.WriteRune(r)
	t.currState = tokenLeftParen
	t.parenDepth.increment(t.currIndex)
	return
}

func (t *tokenizer) handleRightParen(r rune) (err error) {
	// can't allow lone decimal point to be followed by a right parenthesis
	if t.currState == tokenLoneDecimalPoint {
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
		}
	}
	// can't allow unmatched parentheses
	if t.parenDepth.current() == 0 {
		return SyntaxError{
			Message:  fmt.Sprintf(errUnmatchedRightParen, t.currIndex),
			Token:    ")",
			Position: t.currIndex,
		}
	}
	// can't allow empty parentheses
	if t.currState == tokenLeftParen {
		return SyntaxError{
			Message:  fmt.Sprintf(errEmptyParen, t.currIndex),
			Token:    ")",
			Position: t.currIndex,
		}
	}
	// can't allow unfinished operations
	if t.currState&(tokenLeftUnaryOp|tokenBinaryOp) != 0 {
		return SyntaxError{
			Message:  fmt.Sprintf(errNoRightOperand, t.currSymbol.String(), t.currIndex-1),
			Token:    t.currSymbol.String(),
			Position: t.currIndex - 1,
		}
	}

	t.commitCurrentState()
	t.currSymbol.WriteRune(r)
	t.currState = tokenRightParen
	return
}

func (t *tokenizer) handleOperator(r rune) (err error) {
	// can't allow lone decimal point to be followed by operator
	if t.currState == tokenLoneDecimalPoint {
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
		}
	}

	// handle sign
	if (r == '+' || r == '-') && t.currState&(tokenLeftParen|tokenBinaryOp|tokenLeftUnaryOp) != 0 {
		t.commitCurrentState()
		t.currState = tokenLeftUnaryOp
		t.currSymbol.WriteRune(r)
		return
	}

	// at this point, operators should require a left operand.
	if t.currState&(tokenLeftParen|tokenLeftUnaryOp|tokenBinaryOp) != 0 {
		return SyntaxError{
			Message:  fmt.Sprintf(errNoLeftOperand, string(r), t.currIndex),
			Token:    string(r),
			Position: t.currIndex,
		}
	}

	t.commitCurrentState()
	t.currSymbol.WriteRune(r)
	if r == '!' {
		t.currState = tokenRightUnaryOp
	} else {
		t.currState = tokenBinaryOp
	}
	return
}

func (t *tokenizer) commitCurrentState() {
	x := t.currSymbol.String()

	switch t.currState {
	case tokenInteger, tokenDecimal:
		t.appendToken(NewNumber(x))
	case tokenLeftParen:
		t.appendToken(LeftParen{})
	case tokenRightParen:
		t.appendToken(RightParen{})
	case tokenLeftUnaryOp:
		switch x {
		case "+":
			t.appendToken(NewOperator(Plus))
		case "-":
			t.appendToken(NewOperator(Minus))
		}
	case tokenBinaryOp, tokenRightUnaryOp:
		switch x {
		case "+":
			t.appendToken(NewOperator(Addition))
		case "-":
			t.appendToken(NewOperator(Subtraction))
		case "*":
			t.appendToken(NewOperator(Multiplication))
		case "/":
			t.appendToken(NewOperator(Division))
		case "!":
			t.appendToken(NewOperator(Factorial))
		}
	}
	t.currSymbol.Reset()
}

func (t *tokenizer) validateFinalState() (err error) {
	switch t.currState {
	// can't allow lone decimal point as final state
	case tokenLoneDecimalPoint:
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
		}
	case tokenLeftUnaryOp, tokenBinaryOp:
		op := t.currSymbol.String()
		return SyntaxError{
			Message:  fmt.Sprintf(errNoRightOperand, op, t.currIndex),
			Token:    op,
			Position: t.currIndex - 1,
		}
	}

	// check paren depth
	if t.parenDepth.current() > 0 {
		n := len(t.parenDepth.stack)
		curr := t.parenDepth.current()

		var idx int
		for i := n - 1; i >= 0; i-- {
			if t.parenDepth.stack[i][1] == curr {
				idx = t.parenDepth.stack[i][0]
				break
			}
		}

		return SyntaxError{
			Message:  fmt.Sprintf(errUnmatchedLeftParen, idx),
			Token:    "(",
			Position: idx,
		}
	}

	return
}

func (t *tokenizer) appendToken(_t Token) {
	t.tokens = append(t.tokens, _t)
}

func (t *tokenizer) initialize(expr string) {
	t.reset()

	buf := bytes.NewBufferString(expr)
	t.sc = bufio.NewScanner(buf)
	t.sc.Split(bufio.ScanRunes)
}

func (t *tokenizer) reset() {
	t.currState = tokenNothing
	t.currIndex = 0
	t.parenDepth = new(depthStack)
	t.currSymbol.Reset()
}