package pool

import (
	"github.com/ckcfcc/basis/mathx"
	"sync"
)

var bm *BytesMgr = newBytesMgr()

// POT:Power Of Two
const minPOT = 5  // 32 Byte
const maxPOT = 16 // 64 MByte
const minItemNum = 64
const minPOTSize = 2 << (minPOT - 1)
const maxPOTSize = 2 << (maxPOT - 1)

var itemNum uint32 = 64
var bytesMgr = &BytesMgr{}

func SetPoolItemNum(num uint32) {
	if num > minItemNum {
		itemNum = num
	}
}

func GetBytes(size int) []byte {
	return bm.Get(size)
}

func PutBytes(b []byte) {
	bm.Put(b)
}

type BytesMgr struct {
	sync.RWMutex
	_pools map[uint32]*BytePool
}

func newBytesMgr() (bm *BytesMgr) {
	bm = &BytesMgr{}
	bm._pools = make(map[uint32]*BytePool)
	for i := uint32(minPOT); i <= maxPOT; i++ {
		size := uint32(2 << (i - 1))
		bm._pools[size] = NewBytePool(size, itemNum)
	}
	return
}

func (bm *BytesMgr) Get(size int) (b []byte) {
	potSize := uint32(mathx.MinPowerOf2(size))

	if potSize < minPOTSize {
		potSize = minPOTSize
	}

	if potSize > maxPOTSize {
		return make([]byte, potSize)
	}

	return bm._pools[potSize].Get()
}

func (bm *BytesMgr) Put(b []byte) {
	potSize := uint32(len(b))

	if !mathx.IsPowerOf2(int(potSize)) {
		panic("buffer isn't power of 2!")
	}

	if potSize < minPOTSize {
		potSize = minPOTSize
	}

	if potSize > maxPOTSize {
		return
	}

	bm._pools[potSize].Put(b)
}
