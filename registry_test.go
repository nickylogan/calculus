package calculus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatorRegistry_GetOperator(t *testing.T) {
	reg := defaultOperators
	op, _ := reg.GetOperator('+', IsUnaryOp, IsRightAssocOp)
	assert.Equal(t, NewOperator(Plus), op)
	op, _ = reg.GetOperator('+', IsBinaryOp)
	assert.Equal(t, NewOperator(Addition), op)
	op, _ = reg.GetOperator('-', IsUnaryOp, IsRightAssocOp)
	assert.Equal(t, NewOperator(Minus), op)
	op, _ = reg.GetOperator('-', IsBinaryOp)
	assert.Equal(t, NewOperator(Subtraction), op)
	op, _ = reg.GetOperator('*', IsBinaryOp)
	assert.Equal(t, NewOperator(Multiplication), op)
	op, _ = reg.GetOperator('/', IsBinaryOp)
	assert.Equal(t, NewOperator(Division), op)
	op, _ = reg.GetOperator('^', IsBinaryOp)
	assert.Equal(t, NewOperator(Power), op)
	op, _ = reg.GetOperator('!', IsUnaryOp, IsLeftAssocOp)
	assert.Equal(t, NewOperator(Factorial), op)
}
