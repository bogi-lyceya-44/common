package set_test

import (
	"testing"

	"github.com/bogi-lyceya-44/common/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lhs  set.Set[int]
		rhs  set.Set[int]
		want bool
	}{
		{
			name: "equal sets of ints",
			lhs:  set.New(1, 2, 3, 4, 5),
			rhs:  set.New(5, 4, 3, 2, 1),
			want: true,
		},
		{
			name: "non-equal sets of ints",
			lhs:  set.New(1, 5, 3),
			rhs:  set.New(3, 1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				assert.Equal(t, tt.want, tt.lhs.Equal(tt.rhs))
			},
		)
	}
}

func TestSubstitute(t *testing.T) {
	t.Parallel()

	lhs := set.New("abacaba", "babababa")
	rhs := set.New("babababa", "aaaaaa")
	want := set.New("abacaba")

	assert.True(t, want.Equal(lhs.Substitute(rhs)))
}
