package mathx

func MinPow2(val int) int {
	if val == 0 {
		return val
	}

	// 原始值:00001001
	// 右移值:00000100
	// 取或值:00001101
	// 加一值:00001110
	val |= val >> 1
	val += 1
	return val
}
