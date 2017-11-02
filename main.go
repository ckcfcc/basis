package main

import (
	"github.com/ckcfcc/basis/logx"
	"math/rand"
	"time"
)

func main() {
	count := 0
	rvArr := ([]int)(nil)

	for {
		newRv := 5*1024 + rand.Intn(5*1024)
		//log.Printf("第%04d次收包:%d字节\n", count+1, newRv)
		logx.Infof("20秒内平均流量:%f\n", getRate(&rvArr, 20, count, newRv))
		count++
		time.Sleep(time.Millisecond * 500)

		// newRv := 5*1024 + rand.Intn(5*1024)
		// log.Printf("第%04d次收包:%d字节\n", count+1, newRv)
		// if len(rvArr) > 20 {
		// 	rvArr[count%20] = newRv
		// } else {
		// 	rvArr = append(rvArr, newRv)
		// }

		// all := 0
		// rvArrLen := len(rvArr)
		// for _, v := range rvArr {
		// 	all += v
		// }
		// per := float32(all) / float32(rvArrLen)

		// log.Printf("20秒内平均流量:%f\n", per)
		// count++
		// time.Sleep(time.Millisecond * 100
	}
}

func getRate(slots *[]int, secAvg, curIdx, curVal int) (res float32) {
	if *slots == nil {
		*slots = make([]int, 0, secAvg)
	}

	if curIdx+1 > secAvg {
		(*slots)[curIdx%secAvg] = curVal
	} else {
		*slots = append(*slots, curVal)
	}

	tmpAll := 0
	for _, v := range *slots {
		tmpAll += v
	}

	res = float32(tmpAll) / float32(len(*slots))
	return
}
