package yamp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_depthStack_increment(t *testing.T) {
	type fields struct {
		stack []bracketDepth
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStack []bracketDepth
	}{
		{
			name:      "an empty stack should be appended with {idx, 1}",
			fields:    fields{nil},
			args:      args{0},
			wantStack: []bracketDepth{{0, 1, LeftParen}},
		},
		{
			name:      "a filled stack should be appended with {idx, top+1}",
			fields:    fields{[]bracketDepth{{0, 1, LeftParen}}},
			args:      args{3},
			wantStack: []bracketDepth{{0, 1, LeftParen}, {3, 2, LeftParen}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &bracketStack{
				stack: tt.fields.stack,
			}
			d.increment(tt.args.idx, LeftParen)
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}

func Test_depthStack_decrement(t *testing.T) {
	type fields struct {
		stack []bracketDepth
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStack []bracketDepth
	}{
		{
			name:      "an empty stack should be appended with {idx, -1}",
			fields:    fields{nil},
			args:      args{0},
			wantStack: []bracketDepth{{0, -1, RightParen}},
		},
		{
			name:      "a filled stack should be appended with {idx, top-1}",
			fields:    fields{[]bracketDepth{{0, 1, LeftParen}}},
			args:      args{3},
			wantStack: []bracketDepth{{0, 1, LeftParen}, {3, 0, RightParen}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &bracketStack{
				stack: tt.fields.stack,
			}
			d.decrement(tt.args.idx, RightParen)
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}

func Test_depthStack_depth(t *testing.T) {
	type fields struct {
		stack []bracketDepth
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "an empty stack's depth value is 0",
			fields: fields{nil},
		},
		{
			name:   "a filled stack's depth value is the topmost value'",
			fields: fields{[]bracketDepth{{0, 1, LeftParen}}},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &bracketStack{
				stack: tt.fields.stack,
			}
			assert.Equal(t, tt.want, d.depth())
		})
	}
}

func Test_depthStack_clear(t *testing.T) {
	type fields struct {
		stack []bracketDepth
	}
	tests := []struct {
		name      string
		fields    fields
		wantStack []bracketDepth
	}{
		{
			name:   "it should clear the internal stack",
			fields: fields{[]bracketDepth{{0, 1, LeftParen}, {1, 2, LeftParen}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &bracketStack{
				stack: tt.fields.stack,
			}
			d.clear()
			assert.Equal(t, tt.wantStack, d.stack)
		})
	}
}
