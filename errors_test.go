package calculus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntaxError_Error(t *testing.T) {
	type fields struct {
		Message  string
		Token    string
		Position int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"it shall return the internal message", fields{Message: "abc"}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SyntaxError{
				Message:  tt.fields.Message,
				Token:    tt.fields.Token,
				Position: tt.fields.Position,
			}
			assert.Equal(t, tt.want, s.Error())
		})
	}
}
