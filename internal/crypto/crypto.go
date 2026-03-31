package crypto

import (
	"math/rand/v2"
	"slices"
)

const B64CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

func lookupTable(t []byte) [256]int {
	var l [256]int
	for i := range l {
		l[i] = -1
	}

	for i, char := range t {
		l[char] = i
	}

	return l
}

// Basic substitution cipher over plaintext by key, uses B64Charset as base
func RotCipher(pt []byte, key []byte, base []byte) []byte {
	if len(key) != len(base) {
		return nil
	}

	ct := make([]byte, len(pt))
	lt := lookupTable(base)

	for i, b := range pt {
		idx := lt[b]
		if idx == -1 {
			return nil
		}
		ct[i] = key[idx]
	}

	return ct
}

// Shuffles a in place
func Shuffle[E any](a []E) {
	slices.SortFunc(a, func(a E, b E) int {
		return rand.IntN(3) - 1
	})
}

// Shuffles a and b, preserving the mapping between the 2 of them
func ShuffleWith[E any, Z any](a []E, b []Z) {
	s1 := rand.Uint64()
	s2 := rand.Uint64()

	p1 := rand.NewPCG(s1, s2)
	p2 := rand.NewPCG(s1, s2)

	r1 := rand.New(p1)
	r2 := rand.New(p2)

	slices.SortFunc(a, func(a E, b E) int {
		return r1.IntN(3) - 1
	})

	slices.SortFunc(b, func(a Z, b Z) int {
		return r2.IntN(3) - 1
	})
}
