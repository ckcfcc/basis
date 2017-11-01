package pool

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_BytePool(t *testing.T) {
	wg := &sync.WaitGroup{}
	bp := NewBytePool(64, 64)
	rand.Seed(time.Now().Unix())

	use := func(useTime int) {
		for i := 0; i < 100; i++ {
			b := bp.Get()
			// use b
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(useTime)))
			bp.Put(b)
		}
	}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			//log.Printf("start go func %d\n", idx)
			use(10)
		}(i)
	}

	wg.Wait()
	t.Log("BytePool ok")
}
