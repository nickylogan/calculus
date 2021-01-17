package yamp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenizer(t *testing.T) {
	tt := NewTokenizer()
	assert.NotNil(t, tt)
}

func Test_tokenizer_Tokenize(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name       string
		args       args
		wantTokens []Token
		wantErr    error
	}{
		{
			name:       "expr #1",
			args:       args{expr: "56"},
			wantTokens: []Token{NewNumber("56")},
			wantErr:    nil,
		},
		{
			name: "expr #2",
			args: args{expr: "5+6"},
			wantTokens: []Token{
				NewNumber("5"),
				NewOperator(Addition),
				NewNumber("6"),
			},
		},
		{
			name: "expr #3",
			args: args{expr: "(10 +- 7.5) / ((-5)7)"},
			wantTokens: []Token{
				LeftParen,
				NewNumber("10"),
				NewOperator(Addition),
				NewOperator(Minus),
				NewNumber("7.5"),
				RightParen,
				NewOperator(Division),
				LeftParen,
				LeftParen,
				NewOperator(Minus),
				NewNumber("5"),
				RightParen,
				NewOperator(Multiplication),
				NewNumber("7"),
				RightParen,
			},
		},
		{
			name:       "expr #4",
			args:       args{expr: "5###"},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errUnknownSymbol, "#", 1),
				Token:    "#",
				Position: 1,
			},
		},
		{
			name:       "expr #5",
			args:       args{expr: ".+"},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, 0),
				Token:    ".",
				Position: 0,
			},
		},
		{
			name:       "expr #6",
			args:       args{expr: ".("},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, 0),
				Token:    ".",
				Position: 0,
			},
		},
		{
			name:       "expr #7",
			args:       args{expr: "(5*)"},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoRightOperand, "*", 2),
				Token:    "*",
				Position: 2,
			},
		},
		{
			name:       "expr #8",
			args:       args{expr: "5*/6"},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "/", 2),
				Token:    "/",
				Position: 2,
			},
		},
		{
			name:       "expr #9",
			args:       args{expr: "1..5"},
			wantTokens: nil,
			wantErr: &SyntaxError{
				Message:  errMultipleDecimal,
				Token:    ".",
				Position: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTokenizer()
			gotTokens, err := tr.Tokenize(tt.args.expr)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}

func Test_tokenizer_handleDigit(t *testing.T) {
	type fields struct {
		tokens     []Token
		currState  int
		currSymbol string
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantState      int
		wantCurrSymbol string
		wantTokens     []Token
		wantErr        error
	}{
		{
			name:           "a digit should append an existing integer",
			fields:         fields{currState: tokenInteger, currSymbol: "123"},
			args:           args{r: '4'},
			wantState:      tokenInteger,
			wantCurrSymbol: "1234",
			wantErr:        nil,
		},
		{
			name:           "a digit should append an existing decimal",
			fields:         fields{currState: tokenDecimal, currSymbol: "1.5"},
			args:           args{r: '7'},
			wantState:      tokenDecimal,
			wantCurrSymbol: "1.57",
			wantErr:        nil,
		},
		{
			name:           "a digit should convert a decimal point into a decimal",
			fields:         fields{currState: tokenDecimalPoint, currSymbol: "."},
			args:           args{r: '5'},
			wantState:      tokenDecimal,
			wantCurrSymbol: ".5",
			wantErr:        nil,
		},
		{
			name:           "a digit at the start should set the current state into an integer",
			fields:         fields{currState: tokenNothing, currSymbol: ""},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
		},
		{
			name:           "a * operator should be inserted between a ')' and a digit",
			fields:         fields{currState: tokenRightParen, currSymbol: ")"},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
			wantTokens:     []Token{RightParen, NewOperator(Multiplication)},
		},
		{
			name:           "a * operator should be inserted between a right unary op and a digit",
			fields:         fields{currState: tokenRightUnaryOp, currSymbol: "!"},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Factorial), NewOperator(Multiplication)},
		},
		{
			name:           "a binary operator should be committed before adding a digit",
			fields:         fields{currState: tokenBinaryOp, currSymbol: "/"},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Division)},
		},
		{
			name:           "a '(' should be committed before adding a digit",
			fields:         fields{currState: tokenLeftParen, currSymbol: "("},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
			wantTokens:     []Token{LeftParen},
		},
		{
			name:           "a left unary operator should be committed before adding a digit",
			fields:         fields{currState: tokenLeftUnaryOp, currSymbol: "-"},
			args:           args{r: '5'},
			wantState:      tokenInteger,
			wantCurrSymbol: "5",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Minus)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)
			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				tokens:     tt.fields.tokens,
				currState:  tt.fields.currState,
				currSymbol: sb,
			}

			err := tr.handleDigit(tt.args.r)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantState, tr.currState)
			assert.Equal(t, tt.wantCurrSymbol, tr.currSymbol.String())
			assert.Equal(t, tt.wantTokens, tr.tokens)
		})
	}
}

func Test_tokenizer_handleDecimalPoint(t *testing.T) {
	type fields struct {
		tokens     []Token
		currState  int
		currSymbol string
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantState      int
		wantCurrSymbol string
		wantTokens     []Token
		wantErr        error
	}{
		{
			name:           "a decimal point should convert an integer into a decimal",
			fields:         fields{currState: tokenInteger, currSymbol: "123"},
			args:           args{r: '.'},
			wantState:      tokenDecimal,
			wantCurrSymbol: "123.",
			wantErr:        nil,
		},
		{
			name:           "multiple decimal points in a decimal should not be allowed",
			fields:         fields{currState: tokenDecimal, currSymbol: "1.5"},
			args:           args{r: '.'},
			wantState:      tokenDecimal,
			wantCurrSymbol: "1.5",
			wantErr: SyntaxError{
				Message:  errMultipleDecimal,
				Token:    ".",
				Position: 0,
			},
		},
		{
			name:           "a decimal point should not be followed with another decimal point",
			fields:         fields{currState: tokenDecimalPoint, currSymbol: "."},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr: SyntaxError{
				Message:  errMultipleDecimal,
				Token:    ".",
				Position: 0,
			},
		},
		{
			name:           "a decimal point at the start should set the current state into a decimal point",
			fields:         fields{currState: tokenNothing, currSymbol: ""},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
		},
		{
			name:           "a * operator should be inserted between a ')' and a decimal point",
			fields:         fields{currState: tokenRightParen, currSymbol: ")"},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
			wantTokens:     []Token{RightParen, NewOperator(Multiplication)},
		},
		{
			name:           "a * operator should be inserted between a right unary op and a decimal point",
			fields:         fields{currState: tokenRightUnaryOp, currSymbol: "!"},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Factorial), NewOperator(Multiplication)},
		},
		{
			name:           "a binary operator should be committed before adding a decimal point",
			fields:         fields{currState: tokenBinaryOp, currSymbol: "/"},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Division)},
		},
		{
			name:           "a '(' should be committed before adding a decimal point",
			fields:         fields{currState: tokenLeftParen, currSymbol: "("},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
			wantTokens:     []Token{LeftParen},
		},
		{
			name:           "a left unary operator should be committed before adding a digit",
			fields:         fields{currState: tokenLeftUnaryOp, currSymbol: "-"},
			args:           args{r: '.'},
			wantState:      tokenDecimalPoint,
			wantCurrSymbol: ".",
			wantErr:        nil,
			wantTokens:     []Token{NewOperator(Minus)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)
			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				tokens:     tt.fields.tokens,
				currState:  tt.fields.currState,
				currSymbol: sb,
			}

			err := tr.handleDecimalPoint(tt.args.r)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantState, tr.currState)
			assert.Equal(t, tt.wantCurrSymbol, tr.currSymbol.String())
			assert.Equal(t, tt.wantTokens, tr.tokens)
		})
	}
}

func Test_tokenizer_handleLeftParen(t *testing.T) {
	type fields struct {
		tokens     []Token
		currState  int
		currSymbol string
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantState          int
		wantCurrSymbol     string
		wantCurrParenDepth int
		wantTokens         []Token
		wantErr            error
	}{
		{
			name: "a left paren is allowed at the start of an expression",
			fields: fields{
				currState: tokenNothing,
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         nil,
			wantErr:            nil,
		},
		{
			name: "a * operator should be inserted between an integer and a left paren",
			fields: fields{
				currState:  tokenInteger,
				currSymbol: "123",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{NewNumber("123"), NewOperator(Multiplication)},
			wantErr:            nil,
		},
		{
			name: "a single decimal point cannot be followed by a left paren",
			fields: fields{
				currState:  tokenDecimalPoint,
				currSymbol: ".",
			},
			args:               args{r: '('},
			wantState:          tokenDecimalPoint,
			wantCurrSymbol:     ".",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a * operator should be inserted between a decimal and a left paren",
			fields: fields{
				currState:  tokenDecimal,
				currSymbol: "1.23",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{NewNumber("1.23"), NewOperator(Multiplication)},
			wantErr:            nil,
		},
		{
			name: "a left paren should increase the depth stack",
			fields: fields{
				currState:  tokenLeftParen,
				currSymbol: "(",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{LeftParen},
			wantErr:            nil,
		},
		{
			name: "a * operator should be inserted between a right paren and a left paren",
			fields: fields{
				currState:  tokenRightParen,
				currSymbol: ")",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{RightParen, NewOperator(Multiplication)},
			wantErr:            nil,
		},
		{
			name: "a binary operator can be followed by a left paren",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "/",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{NewOperator(Division)},
			wantErr:            nil,
		},
		{
			name: "a left unary operator can be followed by a left paren",
			fields: fields{
				currState:  tokenLeftUnaryOp,
				currSymbol: "+",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{NewOperator(Plus)},
			wantErr:            nil,
		},
		{
			name: "a * operator should be inserted between a right unary operator and a left paren",
			fields: fields{
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
			},
			args:               args{r: '('},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{NewOperator(Factorial), NewOperator(Multiplication)},
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)
			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				tokens:     tt.fields.tokens,
				currState:  tt.fields.currState,
				currSymbol: sb,
			}

			err := tr.handleLeftParen(tt.args.r)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantState, tr.currState)
			assert.Equal(t, tt.wantCurrSymbol, tr.currSymbol.String())
			assert.Equal(t, tt.wantTokens, tr.tokens)
			assert.Equal(t, tt.wantCurrParenDepth, tr.parenDepth.depth())
		})
	}
}

func Test_tokenizer_handleRightParen(t *testing.T) {
	type fields struct {
		tokens     []Token
		currState  int
		currSymbol string
		parenDepth bracketStack
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantState          int
		wantCurrSymbol     string
		wantCurrParenDepth int
		wantTokens         []Token
		wantErr            error
	}{
		{
			name: "a right paren can follow an integer, given the paren stack > 0",
			fields: fields{
				tokens:     nil,
				currState:  tokenInteger,
				currSymbol: "1",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenRightParen,
			wantCurrSymbol:     ")",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1")},
			wantErr:            nil,
		},
		{
			name: "a right paren can succeed a decimal, given the paren stack > 0",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimal,
				currSymbol: "1.5",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenRightParen,
			wantCurrSymbol:     ")",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1.5")},
			wantErr:            nil,
		},
		{
			name: "a right paren cannot exist if the paren stack is empty",
			fields: fields{
				tokens:     nil,
				currState:  tokenInteger,
				currSymbol: "1",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: ')'},
			wantState:          tokenInteger,
			wantCurrSymbol:     "1",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errUnmatchedRightParen, 0),
				Token:    ")",
				Position: 0,
			},
		},
		{
			name: "a right paren cannot succeed a decimal point",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimalPoint,
				currSymbol: ".",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenDecimalPoint,
			wantCurrSymbol:     ".",
			wantCurrParenDepth: 1,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a right paren must not be directly preceded by a left paren",
			fields: fields{
				tokens:     []Token{LeftParen},
				currState:  tokenLeftParen,
				currSymbol: "(",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{LeftParen},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errEmptyParen, 0),
				Token:    ")",
				Position: -1,
			},
		},
		{
			name: "a right paren can follow a right paren, given parent stack > 0",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightParen,
				currSymbol: ")",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenRightParen,
			wantCurrSymbol:     ")",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{RightParen},
			wantErr:            nil,
		},
		{
			name: "a right paren must not succeed an unfinished binary operation",
			fields: fields{
				tokens:     []Token{LeftParen, NewNumber("1"), NewOperator(Multiplication)},
				currState:  tokenBinaryOp,
				currSymbol: "*",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{LeftParen, NewNumber("1"), NewOperator(Multiplication)},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoRightOperand, "*", -1),
				Token:    "*",
				Position: -1,
			},
		},
		{
			name: "a right paren must not succeed an unfinished left unary operation",
			fields: fields{
				tokens:     []Token{LeftParen, NewOperator(Minus)},
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{LeftParen, NewOperator(Minus)},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoRightOperand, "-", -1),
				Token:    "-",
				Position: -1,
			},
		},
		{
			name: "a right paren can follow a right unary operator, given paren stack > 0",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: ')'},
			wantState:          tokenRightParen,
			wantCurrSymbol:     ")",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Factorial)},
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)

			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				tokens:     tt.fields.tokens,
				currState:  tt.fields.currState,
				currSymbol: sb,
				parenDepth: tt.fields.parenDepth,
			}

			err := tr.handleRightParen(tt.args.r)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantState, tr.currState)
			assert.Equal(t, tt.wantCurrSymbol, tr.currSymbol.String())
			assert.Equal(t, tt.wantTokens, tr.tokens)
			assert.Equal(t, tt.wantCurrParenDepth, tr.parenDepth.depth())
		})
	}
}

func Test_tokenizer_handleOperator(t *testing.T) {
	type fields struct {
		tokens     []Token
		currState  int
		currSymbol string
		parenDepth bracketStack
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantState          int
		wantCurrSymbol     string
		wantCurrParenDepth int
		wantTokens         []Token
		wantErr            error
	}{
		{
			name: "a left unary operator can be at the beginning of an expression",
			fields: fields{
				tokens:     nil,
				currState:  tokenNothing,
				currSymbol: "",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr:            nil,
		},
		{
			name: "a binary operator must not be at the beginning of an expression",
			fields: fields{
				tokens:     nil,
				currState:  tokenNothing,
				currSymbol: "",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenNothing,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "*", 0),
				Token:    "*",
				Position: 0,
			},
		},
		{
			name: "a right unary operator must not be at the beginning of an expression",
			fields: fields{
				tokens:     nil,
				currState:  tokenNothing,
				currSymbol: "",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenNothing,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "*", 0),
				Token:    "*",
				Position: 0,
			},
		},
		{
			name: "a left unary operator that follows an integer becomes a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenInteger,
				currSymbol: "1",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1")},
			wantErr:            nil,
		},
		{
			name: "a binary operator can follow an integer",
			fields: fields{
				tokens:     nil,
				currState:  tokenInteger,
				currSymbol: "1",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1")},
			wantErr:            nil,
		},
		{
			name: "a right unary operator can follow an integer",
			fields: fields{
				tokens:     nil,
				currState:  tokenInteger,
				currSymbol: "1",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenRightUnaryOp,
			wantCurrSymbol:     "!",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1")},
			wantErr:            nil,
		},
		{
			name: "a left unary operator must not succeed a decimal point",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimalPoint,
				currSymbol: ".",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenDecimalPoint,
			wantCurrSymbol:     ".",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a binary operator must not succeed a decimal point",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimalPoint,
				currSymbol: ".",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenDecimalPoint,
			wantCurrSymbol:     ".",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a right unary operator must not succeed a decimal point",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimalPoint,
				currSymbol: ".",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenDecimalPoint,
			wantCurrSymbol:     ".",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a left unary operator that follows a decimal number becomes a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimal,
				currSymbol: "1.5",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1.5")},
			wantErr:            nil,
		},
		{
			name: "a binary operator can follow a decimal number",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimal,
				currSymbol: "1.5",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1.5")},
			wantErr:            nil,
		},
		{
			name: "a right unary operator can follow a decimal number",
			fields: fields{
				tokens:     nil,
				currState:  tokenDecimal,
				currSymbol: "1.5",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenRightUnaryOp,
			wantCurrSymbol:     "!",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewNumber("1.5")},
			wantErr:            nil,
		},
		{
			name: "a left unary operator can follow a left paren",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftParen,
				currSymbol: "(",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: '-'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 1,
			wantTokens:         []Token{LeftParen},
			wantErr:            nil,
		},
		{
			name: "a binary operator must not follow a left paren",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftParen,
				currSymbol: "(",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: '*'},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "*", 0),
				Token:    "*",
				Position: 0,
			},
		},
		{
			name: "a right unary operator must not follow a left paren",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftParen,
				currSymbol: "(",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			args:               args{r: '!'},
			wantState:          tokenLeftParen,
			wantCurrSymbol:     "(",
			wantCurrParenDepth: 1,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "!", 0),
				Token:    "!",
				Position: 0,
			},
		},
		{
			name: "a left unary operator following a right paren becomes a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightParen,
				currSymbol: ")",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{RightParen},
			wantErr:            nil,
		},
		{
			name: "a binary operator can follow a right paren",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightParen,
				currSymbol: ")",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{RightParen},
			wantErr:            nil,
		},
		{
			name: "a right unary operator can follow a right paren",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightParen,
				currSymbol: ")",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenRightUnaryOp,
			wantCurrSymbol:     "!",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{RightParen},
			wantErr:            nil,
		},
		{
			name: "a left unary operator can follow a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenBinaryOp,
				currSymbol: "*",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '-'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Multiplication)},
			wantErr:            nil,
		},
		{
			name: "a binary operator must not follow another binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenBinaryOp,
				currSymbol: "*",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '/'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "/", 0),
				Token:    "/",
				Position: 0,
			},
		},
		{
			name: "a right unary operator must not follow a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenBinaryOp,
				currSymbol: "*",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "!", 0),
				Token:    "!",
				Position: 0,
			},
		},
		{
			name: "a left unary operator can follow another left unary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '+'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "+",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Minus)},
			wantErr:            nil,
		},
		{
			name: "a binary operator must not follow a left unary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "*", 0),
				Token:    "*",
				Position: 0,
			},
		},
		{
			name: "a right unary operator must not follow a left unary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenLeftUnaryOp,
			wantCurrSymbol:     "-",
			wantCurrParenDepth: 0,
			wantTokens:         nil,
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoLeftOperand, "!", 0),
				Token:    "!",
				Position: 0,
			},
		},
		{
			name: "a left unary operator that follows a right unary operator becomes a binary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '+'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "+",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Factorial)},
			wantErr:            nil,
		},
		{
			name: "a binary operator can follow a right unary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '*'},
			wantState:          tokenBinaryOp,
			wantCurrSymbol:     "*",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Factorial)},
			wantErr:            nil,
		},
		{
			name: "a right unary operator can follow another right unary operator",
			fields: fields{
				tokens:     nil,
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
				parenDepth: bracketStack{stack: nil},
			},
			args:               args{r: '!'},
			wantState:          tokenRightUnaryOp,
			wantCurrSymbol:     "!",
			wantCurrParenDepth: 0,
			wantTokens:         []Token{NewOperator(Factorial)},
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)

			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				tokens:     tt.fields.tokens,
				currState:  tt.fields.currState,
				currSymbol: sb,
				parenDepth: tt.fields.parenDepth,
			}

			err := tr.handleOperator(tt.args.r)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantState, tr.currState)
			assert.Equal(t, tt.wantCurrSymbol, tr.currSymbol.String())
			assert.Equal(t, tt.wantTokens, tr.tokens)
			assert.Equal(t, tt.wantCurrParenDepth, tr.parenDepth.depth())
		})
	}
}

func Test_tokenizer_commitCurrentState(t *testing.T) {
	type fields struct {
		currState  int
		currSymbol string
	}
	tests := []struct {
		name       string
		fields     fields
		wantTokens []Token
	}{
		{
			name: "an integer state should be committed as a number token",
			fields: fields{
				currState:  tokenInteger,
				currSymbol: "1",
			},
			wantTokens: []Token{NewNumber("1")},
		},
		{
			name: "a decimal state should be committed as a number token",
			fields: fields{
				currState:  tokenDecimal,
				currSymbol: "1.5",
			},
			wantTokens: []Token{NewNumber("1.5")},
		},
		{
			name: "a left paren state should be committed as a left paren",
			fields: fields{
				currState:  tokenLeftParen,
				currSymbol: "(",
			},
			wantTokens: []Token{LeftParen},
		},
		{
			name: "a right paren state should be committed as a right paren",
			fields: fields{
				currState:  tokenRightParen,
				currSymbol: ")",
			},
			wantTokens: []Token{RightParen},
		},
		{
			name: "a plus sign state should be committed as a plus operator",
			fields: fields{
				currState:  tokenLeftUnaryOp,
				currSymbol: "+",
			},
			wantTokens: []Token{NewOperator(Plus)},
		},
		{
			name: "a minus sign state should be committed as a minus operator",
			fields: fields{
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
			},
			wantTokens: []Token{NewOperator(Minus)},
		},
		{
			name: "an add operator state should be committed as an addition operator",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "+",
			},
			wantTokens: []Token{NewOperator(Addition)},
		},
		{
			name: "a subtraction op state should be committed as a subtraction operator",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "-",
			},
			wantTokens: []Token{NewOperator(Subtraction)},
		},
		{
			name: "a multiplication op state should be committed as a multiplication operator",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "*",
			},
			wantTokens: []Token{NewOperator(Multiplication)},
		},
		{
			name: "a division op state should be committed as a division operator",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "/",
			},
			wantTokens: []Token{NewOperator(Division)},
		},
		{
			name: "a factorial op state should be committed as a factorial operator",
			fields: fields{
				currState:  tokenRightUnaryOp,
				currSymbol: "!",
			},
			wantTokens: []Token{NewOperator(Factorial)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)

			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				currState:  tt.fields.currState,
				currSymbol: sb,
			}
			tr.commitCurrentState()

			// always empty after state commit
			assert.Equal(t, tr.currSymbol.Len(), 0)
			assert.Equal(t, tt.wantTokens, tr.tokens)
		})
	}
}

func Test_tokenizer_validateFinalState(t *testing.T) {
	type fields struct {
		currState  int
		currSymbol string
		parenDepth bracketStack
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "ending with a decimal point is an invalid state",
			fields: fields{
				currState:  tokenDecimalPoint,
				currSymbol: ".",
				parenDepth: bracketStack{stack: nil},
			},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errLoneDecimal, -1),
				Token:    ".",
				Position: -1,
			},
		},
		{
			name: "a left unary operation cannot be unfinished",
			fields: fields{
				currState:  tokenLeftUnaryOp,
				currSymbol: "-",
				parenDepth: bracketStack{stack: nil},
			},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoRightOperand, "-", 0),
				Token:    "-",
				Position: 0,
			},
		},
		{
			name: "a binary operation cannot be unfinished",
			fields: fields{
				currState:  tokenBinaryOp,
				currSymbol: "*",
				parenDepth: bracketStack{stack: nil},
			},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errNoRightOperand, "*", 0),
				Token:    "*",
				Position: 0,
			},
		},
		{
			name: "unbalanced parentheses must produce an error",
			fields: fields{
				currState:  tokenLeftParen,
				currSymbol: "(",
				parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
			},
			wantErr: &SyntaxError{
				Message:  fmt.Sprintf(errUnmatchedLeftParen, 0),
				Token:    "(",
				Position: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			sb.WriteString(tt.fields.currSymbol)

			tr := &tokenizer{
				reg:        defaultTokenRegistry,
				currState:  tt.fields.currState,
				currSymbol: sb,
				parenDepth: tt.fields.parenDepth,
			}

			err := tr.validateFinalState()
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_tokenizer_reset(t *testing.T) {
	sb := new(strings.Builder)
	sb.WriteString("abcde")

	tt := &tokenizer{
		currState:  tokenRightParen,
		currIndex:  10,
		parenDepth: bracketStack{stack: []bracketDepth{{0, 1, LeftParen}}},
		currSymbol: sb,
	}
	tt.reset()

	assert.Equal(t, tokenNothing, tt.currState)
	assert.Equal(t, 0, tt.currIndex)
	assert.Equal(t, 0, tt.parenDepth.depth())
	assert.Equal(t, 0, tt.currSymbol.Len())
}
