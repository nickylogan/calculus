package calculus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeftParen_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"a LeftParen's string representation should be '('", "("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LeftParen{}
			assert.Equal(t, tt.want, l.String())
		})
	}
}

func TestRightParen_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"a RightParen's string representation should be ')'", ")"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RightParen{}
			assert.Equal(t, tt.want, r.String())
		})
	}
}

func Test_depthStack_increment(t *testing.T) {
	type fields struct {
		stack [][2]int
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStack [][2]int
	}{
		{"an empty stack should be appended with {idx, 1}", fields{nil}, args{0}, [][2]int{{0, 1}}},
		{"a filled stack should be appended with {idx, top+1}", fields{[][2]int{{0, 1}}}, args{3}, [][2]int{{0, 1}, {3, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &depthStack{
				stack: tt.fields.stack,
			}
			d.increment(tt.args.idx)
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}

func Test_depthStack_decrement(t *testing.T) {
	type fields struct {
		stack [][2]int
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStack [][2]int
	}{
		{"an empty stack should be appended with {idx, -1}", fields{nil}, args{0}, [][2]int{{0, -1}}},
		{"a filled stack should be appended with {idx, top-1}", fields{[][2]int{{0, 1}}}, args{3}, [][2]int{{0, 1}, {3, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &depthStack{
				stack: tt.fields.stack,
			}
			d.decrement(tt.args.idx)
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}

func Test_depthStack_current(t *testing.T) {
	type fields struct {
		stack [][2]int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"an empty stack's current value is 0", fields{nil}, 0},
		{"a filled stack's current value is the topmost value'", fields{[][2]int{{0, 1}}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &depthStack{
				stack: tt.fields.stack,
			}
			assert.Equal(t, tt.want, d.current())
		})
	}
}

func Test_depthStack_clear(t *testing.T) {
	type fields struct {
		stack [][2]int
	}
	tests := []struct {
		name      string
		fields    fields
		wantStack [][2]int
	}{
		{"it should clear the internal stack", fields{[][2]int{{0, 1}, {1, 2}}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &depthStack{
				stack: tt.fields.stack,
			}
			d.clear()
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}
