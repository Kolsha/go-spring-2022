//go:build !solution

package genericsum

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"math/cmplx"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[T constraints.Ordered](a []T) {
	slices.Sort(a)
}

func MapsEqual[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func SliceContains[E comparable](s []E, v E) bool {
	for _, el := range s {
		if el == v {
			return true
		}
	}

	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	res := make(chan T)
	go func() {
		open := len(chs)
		for {
			for _, ch := range chs {
				select {
				case v, ok := <-ch:
					if ok {
						res <- v
						continue
					}
					open--
					if open < 1 {
						close(res)
						return
					}
				default:
					continue
				}

			}
		}

	}()
	return res
}

type Numeric interface {
	constraints.Integer | constraints.Complex | constraints.Float
}

func IsHermitianMatrix[T Numeric](m [][]T) bool {
	height := len(m)
	if height < 1 {
		return true
	}
	width := len(m[0])
	if width < 1 {
		return true
	}
	if height != width {
		return false
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			el := m[x][y]
			switch any(el).(type) {
			case complex64:
				oposite, ok := any(m[y][x]).(complex64)
				if !ok {
					return false
				}
				cur, ok := any(m[x][y]).(complex64)
				if !ok {
					return false
				}

				if real(oposite) != real(cur) {
					return false
				}
				if -imag(oposite) != imag(cur) {
					return false
				}

			case complex128:
				oposite, ok := any(m[y][x]).(complex128)
				if !ok {
					return false
				}
				cur, ok := any(m[x][y]).(complex128)
				if !ok {
					return false
				}

				if cmplx.Conj(cur) != oposite {
					return false
				}
			default:
				if m[x][y] != m[y][x] {
					return false
				}
			}
		}
	}

	return true
}
