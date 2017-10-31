package mathx

import ()

func MinPowerOf2(n int) int {
	if n == 0 {
		return 1
	}

	n -= 1
	n |= n >> 16
	n |= n >> 8
	n |= n >> 4
	n |= n >> 2
	n |= n >> 1
	return n + 1
}

func IsPowerOf2(n int) bool {
	return (n & -n) == n
}
