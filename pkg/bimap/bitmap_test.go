package bimap_test

import (
	"testing"

	"github.com/bogi-lyceya-44/common/pkg/bimap"
	"github.com/bogi-lyceya-44/common/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(b *bimap.Bimap[string, int])
		args  string
		want  utils.Pair[int, bool]
	}{
		{
			name: "default",
			setup: func(b *bimap.Bimap[string, int]) {
				b.Put("one", 1)
			},
			args: "one",
			want: utils.Pair[int, bool]{First: 1, Second: true},
		},
		{
			name:  "key not found",
			setup: func(b *bimap.Bimap[string, int]) {},
			args:  "one",
			want:  utils.Pair[int, bool]{First: 0, Second: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			bm := bimap.New[string, int]()
			tt.setup(bm)

			v, ok := bm.Get(tt.args)

			assert.Equal(st, tt.want.First, v)
			assert.Equal(st, tt.want.Second, ok)
		})
	}
}

func TestGetInverse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(b *bimap.Bimap[string, int])
		args  int
		want  utils.Pair[string, bool]
	}{
		{
			name: "default",
			setup: func(b *bimap.Bimap[string, int]) {
				b.Put("one", 1)
			},
			args: 1,
			want: utils.Pair[string, bool]{First: "one", Second: true},
		},
		{
			name:  "key not found",
			setup: func(b *bimap.Bimap[string, int]) {},
			args:  1,
			want:  utils.Pair[string, bool]{First: "", Second: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			bm := bimap.New[string, int]()
			tt.setup(bm)

			v, ok := bm.GetInverse(tt.args)

			assert.Equal(st, tt.want.First, v)
			assert.Equal(st, tt.want.Second, ok)
		})
	}
}

func TestNewFromMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		args           string
		argsInverse    int
		wantGet        utils.Pair[int, bool]
		wantGetInverse utils.Pair[string, bool]
	}{
		{
			name:           "default",
			args:           "one",
			argsInverse:    1,
			wantGet:        utils.Pair[int, bool]{First: 1, Second: true},
			wantGetInverse: utils.Pair[string, bool]{First: "one", Second: true},
		},
		{
			name:           "key not found",
			args:           "lol",
			argsInverse:    52,
			wantGet:        utils.Pair[int, bool]{First: 0, Second: false},
			wantGetInverse: utils.Pair[string, bool]{First: "", Second: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			bm := bimap.NewFromMap(
				map[string]int{"one": 1},
			)

			vGet, okGet := bm.Get(tt.args)
			vGetInverse, okGetInverse := bm.GetInverse(tt.argsInverse)

			assert.Equal(st, tt.wantGet.First, vGet)
			assert.Equal(st, tt.wantGet.Second, okGet)

			assert.Equal(st, tt.wantGetInverse.First, vGetInverse)
			assert.Equal(st, tt.wantGetInverse.Second, okGetInverse)
		})
	}
}
