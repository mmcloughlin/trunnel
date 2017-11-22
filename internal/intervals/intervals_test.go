package intervals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	assert.Equal(t, Interval{lo: 10, hi: 20}, Range(10, 20))
}

func TestBadRange(t *testing.T) {
	assert.Panics(t, func() { Range(4, 3) })
}

func TestSingle(t *testing.T) {
	assert.Equal(t, Interval{lo: 13, hi: 13}, Single(13))
}

func TestIntervalSize(t *testing.T) {
	assert.Equal(t, uint64(1), Single(13).Size())
	assert.Equal(t, uint64(5), Range(13, 17).Size())
}

func TestIntervalString(t *testing.T) {
	assert.Equal(t, "13", Single(13).String())
	assert.Equal(t, "13-17", Range(13, 17).String())
}

func TestSetOverlaps(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		Expect    bool
	}{
		{[]Interval{}, false},
		{[]Interval{Single(1)}, false},
		{[]Interval{Range(10, 100)}, false},
		{[]Interval{Range(10, 100), Range(50, 60)}, true},
		{[]Interval{Range(5, 10), Range(10, 15)}, true},
		{[]Interval{Range(5, 10), Range(11, 15)}, false},
	}
	for _, c := range cases {
		s := Set(c.Intervals)
		t.Run(s.String(), func(t *testing.T) {
			assert.Equal(t, c.Expect, s.Overlaps())
		})
	}
}

func TestSetString(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		Expect    string
	}{
		{
			Intervals: []Interval{},
			Expect:    "",
		},
		{
			Intervals: []Interval{Single(2), Range(4, 5)},
			Expect:    "2,4-5",
		},
		{
			Intervals: []Interval{Single(2), Range(4, 50), Range(30, 300)},
			Expect:    "2,4-50,30-300",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.Expect, Set(c.Intervals).String())
	}
}
