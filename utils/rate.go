package utils

// slots:一个[]int的切片，用于保存所有需要统计的数据 该切片会由函数自动创建
// secAvgNum:计算多少时间内的平均流量
// curVal:当前值
// curIdx:当前数据索引 每次调用时此值应该递增
func GetAvgRate(slots *[]int, secAvgNum, curVal, curIdx int) (res float64) {
	if *slots == nil {
		*slots = make([]int, 0, secAvgNum)
	}

	if curIdx+1 > secAvgNum {
		(*slots)[curIdx%secAvgNum] = curVal
	} else {
		*slots = append(*slots, curVal)
	}

	tmpAll := 0
	for _, v := range *slots {
		tmpAll += v
	}

	res = float64(tmpAll) / float64(len(*slots))
	return
}
