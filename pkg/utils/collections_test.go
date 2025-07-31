package utils_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/bogi-lyceya-44/common/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	got := utils.Map(
		arr,
		func(x int) int {
			return x * 2
		},
	)

	assert.Equal(t, expected, got)
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	type item struct {
		name string
	}

	item1 := &item{name: "lol"}
	item2 := &item{name: "abc"}
	item3 := &item{name: "lol"}
	item4 := &item{name: "cba"}
	arr := []*item{item1, item2, item3, item4}

	expected := map[string][]*item{
		"lol": {item1, item3},
		"abc": {item2},
		"cba": {item4},
	}

	got := utils.GroupBy(
		arr,
		func(x *item) string {
			return x.name
		},
	)

	assert.True(t, reflect.DeepEqual(expected, got))
}

func TestGroupBySingle(t *testing.T) {
	t.Parallel()

	type item struct {
		name string
	}

	item1 := &item{name: "lol"}
	item2 := &item{name: "abc"}
	item3 := &item{name: "cba"}
	arr := []*item{item1, item2, item3}

	expected := map[string][]*item{
		"lol": {item1},
		"abc": {item2},
		"cba": {item3},
	}

	got := utils.GroupBy(
		arr,
		func(x *item) string {
			return x.name
		},
	)

	assert.True(t, reflect.DeepEqual(expected, got))
}

func TestSplit(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}

	expected := [][]int{
		{2, 4},
		{1, 3, 5},
	}

	got := utils.Split(
		arr,
		func(x int) bool {
			return x%2 == 0
		},
	)

	assert.Equal(t, expected, got)
}

func TestFilter(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	expected := []int{1, 3, 5}

	got := utils.Filter(
		arr,
		func(x int) bool {
			return x%2 != 0
		},
	)

	assert.Equal(t, expected, got)
}

func TestBatch(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5, 6, 7}

	expected := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7},
	}

	got := utils.Batch(arr, 3)

	assert.Equal(t, expected, got)
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	arr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7},
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	got := utils.Flatten(arr)

	assert.Equal(t, expected, got)
}

func TestUnique(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 1, 2, 3}
	expected := []int{1, 2, 3}
	got := utils.Unique(arr)

	assert.Equal(t, expected, got)
}

func TestKeys(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"lol": 1,
		"abc": 2,
		"kek": 3,
	}

	expected := []string{"lol", "abc", "kek"}
	got := utils.Keys(m)

	slices.Sort(expected)
	slices.Sort(got)

	assert.Equal(t, expected, got)
}

func TestValues(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"lol": 1,
		"abc": 2,
		"kek": 3,
	}

	expected := []int{1, 2, 3}
	got := utils.Values(m)

	slices.Sort(got)

	assert.Equal(t, expected, got)
}

func TestZip(t *testing.T) {
	t.Parallel()

	lhs := []string{"aaa", "bbb", "ccc", "ddd"}
	rhs := []int{1, 2}

	expected := []utils.Pair[string, int]{
		{"aaa", 1},
		{"bbb", 2},
	}

	got := utils.Zip(lhs, rhs)

	assert.Equal(t, expected, got)
}

func TestReduce(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	expected := 15

	got := utils.Reduce(
		arr,
		func(accumulator int, x int) int {
			return accumulator + x
		},
		0,
	)

	assert.Equal(t, expected, got)
}
