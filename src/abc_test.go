package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	a, b := 1, 2
	if got, want := Add(a, b), 3; got != want {
		t.Errorf("Add(%v, %v) = %v, want: %v", a, b, got, want)
	}

	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{a: 1, b: 2, want: 3, name: "test 1"},
		{a: -2, b: 2, want: 0, name: "test 2"},
		{a: 0, b: 0, want: 0, name: "test 3"},
		{a: 1, b: 2, want: 4, name: "test 4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, Add(tt.a, tt.b), tt.want)
		})
	}
}
