package mathx

import (
	"testing"
)

func Test_MinPow2(t *testing.T) {
	for i := 0; i < 512; i++ {
		t.Logf("[%d]:[%d]\n", i, MinPowerOf2(i))
	}
}

func Test_IsPowerOf2(t *testing.T) {
	t.Logf("%d is power of 2 %v", 2, IsPowerOf2(2))
	t.Logf("%d is power of 2 %v", 4, IsPowerOf2(4))
	t.Logf("%d is power of 2 %v", 8, IsPowerOf2(8))
	t.Logf("%d is power of 2 %v", 10, IsPowerOf2(10))
	t.Logf("%d is power of 2 %v", 16, IsPowerOf2(16))
	t.Logf("%d is power of 2 %v", 32, IsPowerOf2(32))
	t.Logf("%d is power of 2 %v", 64, IsPowerOf2(64))
	t.Logf("%d is power of 2 %v", 256, IsPowerOf2(256))
}
