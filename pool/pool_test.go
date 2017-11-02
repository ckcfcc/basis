package pool

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_BytePoolMgr(t *testing.T) {
	rand.Seed(time.Now().Unix())
	wg := &sync.WaitGroup{}

	// test bytes size array
	tbs := []int{23, 67, 145, 270, 590, 1157, 4934, 10586, 23336}
	tbl := len(tbs)

	v := uint32(128)
	SetPoolItemNum(v)
	if itemNum != v {
		t.Errorf("set pool item num to %d failed! current value:%d", v, itemNum)
	}

	use := func(needSize, needTime int) {
		count := rand.Intn(1024)
		for i := 0; i < count; i++ {
			b := GetBytes(needSize)
			// use b
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(needTime)))
			PutBytes(b)
		}
	}

	for i := 0; i < tbl; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			use(tbs[idx], rand.Intn(512))
		}(i)
	}

	wg.Wait()
	t.Log("BytePoolMgr ok")
}
