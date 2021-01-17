package yamp

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_number_String(t *testing.T) {
	type fields struct {
		symbol string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"it should return the internal symbol", fields{symbol: "5"}, "5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := number{
				symbol: tt.fields.symbol,
			}
			assert.Equal(t, tt.want, n.String())
		})
	}
}

func Test_number_Value(t *testing.T) {
	type fields struct {
		symbol string
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes float64
		wantErr error
	}{
		{
			name:    "it should return the numerical value of its symbol",
			fields:  fields{symbol: "5"},
			wantRes: 5,
			wantErr: nil,
		},
		{
			name:    "it should return an error for a non-numerical symbol",
			fields:  fields{symbol: "a"},
			wantRes: 0,
			wantErr: errors.New("'a' is not a number"),
		},
		{
			name:    "it should return +Inf for extremely large positive numbers",
			fields:  fields{symbol: "1e+100000000"},
			wantRes: math.Inf(1),
			wantErr: nil,
		},
		{
			name:    "it should return -Inf for extremely large negative numbers",
			fields:  fields{symbol: "-1e+100000000"},
			wantRes: math.Inf(-1),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := number{
				symbol: tt.fields.symbol,
			}
			gotRes, err := n.Value()
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestNewNumber(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want Number
	}{
		{"it should create a new number with the given number", args{num: "5"}, number{symbol: "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewNumber(tt.args.num))
		})
	}
}
