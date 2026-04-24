package bitset

import (
	"math/bits"

	"github.com/swonky/set"
)

var _ set.MutableSet[int] = (*BitSet)(nil)

// BitSet is a mutable set of non-negative integers.
type BitSet struct {
	words []uint64
	n     int
}

func (s *BitSet) Contains(item int) bool {
	if item < 0 {
		return false
	}

	word := item >> 6
	if word >= len(s.words) {
		return false
	}

	mask := uint64(1) << uint(item&63)
	return s.words[word]&mask != 0
}

func (s *BitSet) Add(item int) {
	if item < 0 {
		panic("bitset: negative value")
	}

	word := item >> 6
	s.grow(word + 1)

	mask := uint64(1) << uint(item&63)
	if s.words[word]&mask == 0 {
		s.words[word] |= mask
		s.n++
	}
}

func (s *BitSet) Delete(item int) {
	if item < 0 {
		return
	}

	word := item >> 6
	if word >= len(s.words) {
		return
	}

	mask := uint64(1) << uint(item&63)
	if s.words[word]&mask != 0 {
		s.words[word] &^= mask
		s.n--
	}
}

func (s *BitSet) Len() int {
	return s.n
}

func (s *BitSet) Range(yield func(int) bool) {
	for wi, word := range s.words {
		for word != 0 {
			bit := bits.TrailingZeros64(word)
			value := wi<<6 | bit

			if !yield(value) {
				return
			}

			word &^= uint64(1) << uint(bit)
		}
	}
}

func (s *BitSet) grow(n int) {
	if n <= len(s.words) {
		return
	}

	words := make([]uint64, n)
	copy(words, s.words)
	s.words = words
}
