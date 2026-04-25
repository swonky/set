package bitset

import (
	"math/bits"

	"github.com/swonky/set"
)

var _ set.MutableSet[int] = (*BitSet[int])(nil)

type numbers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint16 | ~uint32 | ~uint64
}

// BitSet is a mutable set of non-negative integers.
type BitSet[T numbers] struct {
	words []uint64
	n     int
}

func (s *BitSet[T]) Contains(item int) bool {
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

func (s *BitSet[T]) Add(elem T) {
	if elem < 0 {
		panic("bitset: negative value")
	}

	word := elem >> 6
	s.grow(word + 1)

	mask := uint64(1) << uint(elem&63)
	if s.words[word]&mask == 0 {
		s.words[word] |= mask
		s.n++
	}
}

func (s *BitSet[T]) Delete(item int) {
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

func (s *BitSet[T]) Len() int {
	return s.n
}

func (s *BitSet[T]) Range(yield func(int) bool) {
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

func (s *BitSet[T]) grow(n T) {
	if int(n) <= len(s.words) {
		return
	}

	words := make([]uint64, n)
	copy(words, s.words)
	s.words = words
}
