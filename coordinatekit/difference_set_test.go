package coordinatekit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_difference(t *testing.T) {
	type args struct {
		a Coordinate
		b Coordinate
	}
	tests := []struct {
		name string
		args args
		want []Coordinate
	}{
		{
			name: "four_coordinate",
			args: args{
				a: Coordinate{5, 2, 6, 3},
				b: Coordinate{3, 1, 9, 5},
			},
			want: []Coordinate{
				{3, 1, 5, 5},
				{6, 1, 9, 5},
				{5, 3, 6, 5},
				{5, 1, 6, 2},
			},
		},
		{
			name: "three_coordinate",
			args: args{
				a: Coordinate{5, 2, 6, 3},
				b: Coordinate{5, 1, 9, 5},
			},
			want: []Coordinate{
				{6, 1, 9, 5},
				{5, 3, 6, 5},
				{5, 1, 6, 2},
			},
		},
		{
			name: "two_coordinate",
			args: args{
				a: Coordinate{5, 2, 6, 3},
				b: Coordinate{5, 2, 9, 5},
			},
			want: []Coordinate{
				{6, 2, 9, 5},
				{5, 3, 6, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, difference(tt.args.a, tt.args.b), tt.want)
		})
		// t.Run(tt.name, func(t *testing.T) {
		// 	if got := difference(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
		// 		t.Errorf("difference() = %v, want %v", got, tt.want)
		// 	}
		// })
	}
}
