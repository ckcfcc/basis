package utils

import (
	"math/rand"
	"testing"
	"time"
)

func Test_AvgRate(t *testing.T) {
	rand.Seed(time.Now().Unix())
	rvArr := ([]int)(nil)
	allTime := 10
	for i := 0; i < allTime; i++ {
		newRv := 5*1024 + rand.Intn(5*1024)
		t.Logf("%d秒内平均流量:%f\n", allTime, GetAvgRate(&rvArr, allTime, newRv, i))
		time.Sleep(time.Millisecond * 100)
	}
}
