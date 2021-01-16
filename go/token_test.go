package calculus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDigit(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a digit should be marked as true", args{r: '1'}, true},
		{"a latin alphabet should be marked as false", args{r: 'a'}, false},
		{"a non-alphanumeric ascii symbol should be marked as false", args{r: '$'}, false},
		{"a non-alphanumeric UTF+8 symbol should be marked as false", args{r: '是'}, false},
		{"a whitespace should be marked as false", args{r: ' '}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsDigit(tt.args.r))
		})
	}
}

func TestIsOperator(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"'+' is an operator", args{r: '+'}, true},
		{"'-' is an operator", args{r: '-'}, true},
		{"'*' is an operator", args{r: '*'}, true},
		{"'/' is an operator", args{r: '/'}, true},
		{"'^' is an operator", args{r: '^'}, true},
		{"'!' is an operator", args{r: '!'}, true},
		{"other characters are not an operator", args{r: '?'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsOperator(tt.args.r))
		})
	}
}

func TestIsLeftParen(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a left parenthesis should be marked as true", args{'('}, true},
		{"a right parenthesis should be marked as false", args{')'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsLeftParen(tt.args.r))
		})
	}
}

func TestIsRightParen(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a right parenthesis should be marked as true", args{')'}, true},
		{"a left parenthesis should be marked as false", args{'('}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsRightParen(tt.args.r))
		})
	}
}

func TestIsDecimalPoint(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a period should be marked as true", args{'.'}, true},
		{"a comma should be marked as false", args{','}, false},
		{"other characters should be marked as false", args{'*'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsDecimalPoint(tt.args.r))
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a space should be marked as true", args{' '}, true},
		{"a tab should be marked as true", args{'\t'}, true},
		{"a newline should be marked as true", args{'\n'}, true},
		{"a carriage-return should be marked as true", args{'\r'}, true},
		{"other characters should be marked as false", args{'!'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsWhitespace(tt.args.r))
		})
	}
}
