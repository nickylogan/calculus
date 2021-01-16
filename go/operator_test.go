package calculus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOperator(t *testing.T) {
	type args struct {
		op opType
	}
	tests := []struct {
		name string
		args args
		want Operator
	}{
		{"it should correctly create a new Addition operator", args{op: Addition}, operator{Addition}},
		{"it should correctly create a new Subtraction operator", args{op: Subtraction}, operator{Subtraction}},
		{"it should correctly create a new Multiplication operator", args{op: Multiplication}, operator{Multiplication}},
		{"it should correctly create a new Division operator", args{op: Division}, operator{Division}},
		{"it should correctly create a new Power operator", args{op: Power}, operator{Power}},
		{"it should correctly create a new Plus operator", args{op: Plus}, operator{Plus}},
		{"it should correctly create a new Minus operator", args{op: Minus}, operator{Minus}},
		{"it should correctly create a new Factorial operator", args{op: Factorial}, operator{Factorial}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewOperator(tt.args.op))
		})
	}
}

func Test_operator_String(t *testing.T) {
	type fields struct {
		opType opType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"it should return a '+' for an Addition operator", fields{opType: Addition}, "+"},
		{"it should return a '-' for a Subtraction operator", fields{opType: Subtraction}, "-"},
		{"it should return a '*' for a Multiplication operator", fields{opType: Multiplication}, "*"},
		{"it should return a '/' for a Division operator", fields{opType: Division}, "/"},
		{"it should return a '^' for a Power operator", fields{opType: Power}, "^"},
		{"it should return a '+' for a Plus operator", fields{opType: Plus}, "+"},
		{"it should return a '-' for a Minus operator", fields{opType: Minus}, "-"},
		{"it should return a '!' for a Factorial operator", fields{opType: Factorial}, "!"},
		{"it should return a '<?>' for an unknown operator", fields{opType: -1}, "<?>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operator{
				opType: tt.fields.opType,
			}
			assert.Equal(t, tt.want, o.String())
		})
	}
}

func Test_operator_Type(t *testing.T) {
	type fields struct {
		opType opType
	}
	tests := []struct {
		name   string
		fields fields
		want   opType
	}{
		{"it should return the Addition type for an Addition operator", fields{opType: Addition}, Addition},
		{"it should return the Subtraction type for a Subtraction operator", fields{opType: Subtraction}, Subtraction},
		{"it should return the Multiplication type for a Multiplication operator", fields{opType: Multiplication}, Multiplication},
		{"it should return the Division type for a Division operator", fields{opType: Division}, Division},
		{"it should return the Power type for a Power operator", fields{opType: Power}, Power},
		{"it should return the Plus type for a Plus operator", fields{opType: Plus}, Plus},
		{"it should return the Minus type for a Minus operator", fields{opType: Minus}, Minus},
		{"it should return the Factorial type for a Factorial operator", fields{opType: Factorial}, Factorial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operator{
				opType: tt.fields.opType,
			}
			assert.Equal(t, tt.want, o.Type())
		})
	}
}

func Test_operator_Precedence(t *testing.T) {
	tests := []struct {
		name  string
		type1 opType
		type2 opType
		gt    bool
		eq    bool
	}{
		{
			name:  "addition and subtraction should have the same precedence",
			type1: Addition, type2: Subtraction,
			eq: true,
		},
		{
			name:  "multiplication and division should have the same precedence",
			type1: Multiplication, type2: Division,
			eq: true,
		},
		{
			name:  "multiplication should have a higher precedence than addition",
			type1: Multiplication, type2: Addition,
			gt: true,
		},
		{
			name:  "sign operators should have the same precedence",
			type1: Plus, type2: Minus,
			eq: true,
		},
		{
			name:  "sign operators should have a higher precedence than multiplication",
			type1: Plus, type2: Multiplication,
			gt: true,
		},
		{
			name:  "exponents should have a higher precedence than all basic arithmetic and sign operators",
			type1: Power, type2: Multiplication,
			gt: true,
		},
		{
			name:  "factorial should have the highest precedence",
			type1: Factorial, type2: Power,
			gt: true,
		},
		{
			name:  "unknown operators should have the lowest precedence",
			type1: Addition, type2: -1,
			gt: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op1, op2 := operator{tt.type1}, operator{tt.type2}
			if tt.eq {
				assert.Equal(t, op1.Precedence(), op2.Precedence())
			}
			if tt.gt {
				assert.Greater(t, op1.Precedence(), op2.Precedence())
			}
		})
	}
}

func Test_operator_Associativity(t *testing.T) {
	type fields struct {
		opType opType
	}
	tests := []struct {
		name   string
		fields fields
		want   assoc
	}{
		{"an Addition operator is left associative", fields{opType: Addition}, LeftAssoc},
		{"a Subtraction operator is left associative", fields{opType: Subtraction}, LeftAssoc},
		{"a Multiplication operator is left associative", fields{opType: Multiplication}, LeftAssoc},
		{"a Division operator is left associative", fields{opType: Division}, LeftAssoc},
		{"a Power operator is right associative", fields{opType: Power}, RightAssoc},
		{"a Plus operator is right associative", fields{opType: Plus}, RightAssoc},
		{"a Minus operator is right associative", fields{opType: Minus}, RightAssoc},
		{"a Factorial operator is left associative", fields{opType: Factorial}, LeftAssoc},
		{"an unknown operator is by default left associative", fields{opType: -1}, LeftAssoc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operator{
				opType: tt.fields.opType,
			}
			assert.Equal(t, tt.want, o.Associativity())
		})
	}
}

func Test_operator_Arity(t *testing.T) {
	type fields struct {
		opType opType
	}
	tests := []struct {
		name   string
		fields fields
		want   arity
	}{
		{"Addition operator is a binary operator", fields{opType: Addition}, Binary},
		{"Subtraction operator is a binary operator", fields{opType: Subtraction}, Binary},
		{"Multiplication operator is a binary operator", fields{opType: Multiplication}, Binary},
		{"Division operator is a binary operator", fields{opType: Division}, Binary},
		{"Power operator is a binary operator", fields{opType: Power}, Binary},
		{"Plus operator is a unary operator", fields{opType: Plus}, Unary},
		{"Minus operator is a unary operator", fields{opType: Minus}, Unary},
		{"Factorial operator is a unary operator", fields{opType: Factorial}, Unary},
		{"an unknown operator is by default binary", fields{opType: -1}, Binary},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operator{
				opType: tt.fields.opType,
			}
			assert.Equal(t, tt.want, o.Arity())
		})
	}
}
