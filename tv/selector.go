package tv

import "math/rand"

// Selector selects which vectors to keep from a given list.
type Selector interface {
	SelectVectors([]Vector) []Vector
}

// SelectorFunc implements Selector interface with a plain function.
type SelectorFunc func([]Vector) []Vector

// SelectVectors calls f.
func (f SelectorFunc) SelectVectors(vs []Vector) []Vector {
	return f(vs)
}

// Exhaustive selects all vectors.
var Exhaustive Selector = SelectorFunc(func(vs []Vector) []Vector { return vs })

// RandomSampleSelector selects a random sample of up to n vectors.
func RandomSampleSelector(n int) Selector {
	return SelectorFunc(func(vs []Vector) []Vector {
		m := len(vs)
		if m <= n {
			return vs
		}
		sample := make([]Vector, n)
		pi := rand.Perm(m)
		for i := 0; i < n; i++ {
			sample[i] = vs[pi[i]]
		}
		return sample
	})
}
