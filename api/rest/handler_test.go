package rest

import (
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		a int
		b int
		want int
	}{
		{
			name: "dummy sum",
			a: 1,
			b: 2,
			want: 3,
		},
		{
			name: "negative num",
			a: 1,
			b: -1,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a + tt.b
			if got != tt.want {
				t.Error("Sum test failed")
			}
		})
	}
}