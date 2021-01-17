package yamp

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
	tokenDecimalPoint
	tokenDecimal
	tokenLeftParen
	tokenRightParen
	tokenBinaryOp
	tokenLeftUnaryOp
	tokenRightUnaryOp
)

type tokenizer struct {
	reg        TokenRegistry
	expr       string
	sc         *bufio.Scanner
	tokens     []Token
	currState  int
	currSymbol *strings.Builder
	currIndex  int
	parenDepth bracketStack
}

// NewTokenizer creates a new Tokenizer
func NewTokenizer() Tokenizer {
	return &tokenizer{
		// TODO: add registry validation
		reg:        defaultTokenRegistry,
		currSymbol: new(strings.Builder),
	}
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
		case t.reg.IsDigit(r):
			if err = t.handleDigit(r); err != nil {
				return
			}
		case t.reg.IsDecimalPoint(r):
			if err = t.handleDecimalPoint(r); err != nil {
				return
			}
		case t.reg.IsLeftBracket(r):
			if err = t.handleLeftParen(r); err != nil {
				return
			}
		case t.reg.IsRightBracket(r):
			if err = t.handleRightParen(r); err != nil {
				return
			}
		case t.reg.IsOperator(r):
			if err = t.handleOperator(r); err != nil {
				return
			}
		case t.reg.IsWhitespace(r):
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
	if t.currState == tokenDecimalPoint {
		t.currSymbol.WriteRune(r)
		t.currState = tokenDecimal
		return
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
	if t.currState&(tokenDecimal|tokenDecimalPoint) != 0 {
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
	t.currState = tokenDecimalPoint
	return
}

func (t *tokenizer) handleLeftParen(r rune) (err error) {
	// can't allow lone decimal point to be followed by a left parenthesis
	if t.currState == tokenDecimalPoint {
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
	// TODO: in the future, it may not be only left paren
	t.parenDepth.increment(t.currIndex, LeftParen)
	return
}

func (t *tokenizer) handleRightParen(r rune) (err error) {
	// can't allow unmatched parentheses
	if t.parenDepth.depth() == 0 {
		return SyntaxError{
			Message:  fmt.Sprintf(errUnmatchedRightParen, t.currIndex),
			Token:    ")",
			Position: t.currIndex,
		}
	}
	// can't allow lone decimal point to be followed by a right parenthesis
	if t.currState == tokenDecimalPoint {
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
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
	// TODO: in the future, it may not be only right paren
	t.parenDepth.decrement(t.currIndex, RightParen)
	return
}

func (t *tokenizer) handleOperator(r rune) (err error) {
	// can't allow lone decimal point to be followed by operator
	if t.currState == tokenDecimalPoint {
		return SyntaxError{
			Message:  fmt.Sprintf(errLoneDecimal, t.currIndex-1),
			Token:    ".",
			Position: t.currIndex - 1,
		}
	}

	// handle left unary operators
	_, lunOk := t.reg.operators.GetOperator(r, IsUnaryOp, IsRightAssocOp)
	if lunOk && t.currState&(tokenNothing|tokenLeftParen|tokenBinaryOp|tokenLeftUnaryOp) != 0 {
		t.commitCurrentState()
		t.currState = tokenLeftUnaryOp
		t.currSymbol.WriteRune(r)
		return
	}

	// operator is left unary ONLY, but has no right operands.
	_, binOk := t.reg.operators.GetOperator(r, IsBinaryOp)
	_, runOk := t.reg.operators.GetOperator(r, IsUnaryOp, IsLeftAssocOp)
	if !(binOk || runOk) {
		return SyntaxError{
			Message:  fmt.Sprintf(errNoRightOperand, string(r), t.currIndex),
			Token:    string(r),
			Position: t.currIndex,
		}
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
	if runOk {
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
		t.appendToken(LeftParen)
	case tokenRightParen:
		t.appendToken(RightParen)
	case tokenLeftUnaryOp:
		r := []rune(x)[0]
		op, _ := t.reg.operators.GetOperator(r, IsRightAssocOp, IsUnaryOp)
		t.appendToken(op)
	case tokenBinaryOp:
		r := []rune(x)[0]
		op, _ := t.reg.operators.GetOperator(r, IsBinaryOp)
		t.appendToken(op)
	case tokenRightUnaryOp:
		r := []rune(x)[0]
		op, _ := t.reg.operators.GetOperator(r, IsLeftAssocOp, IsUnaryOp)
		t.appendToken(op)
	}
	t.currSymbol.Reset()
}

func (t *tokenizer) validateFinalState() (err error) {
	switch t.currState {
	// can't allow lone decimal point as final state
	case tokenDecimalPoint:
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
	if t.parenDepth.depth() > 0 {
		n := len(t.parenDepth.stack)
		curr := t.parenDepth.depth()

		var idx int
		for i := n - 1; i >= 0; i-- {
			if t.parenDepth.stack[i].depth == curr {
				idx = t.parenDepth.stack[i].index
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
	t.parenDepth = bracketStack{}
	t.currSymbol.Reset()
}
