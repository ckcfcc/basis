package pool

import (
	"log"
	"sync"
)

type BytePool struct {
	sync.RWMutex
	_readPos    uint32
	_writePos   uint32
	_bytesVec   [][]byte
	_size, _num uint32
}

func NewBytePool(size, num uint32) (bp *BytePool) {
	bp = &BytePool{
		_bytesVec: make([][]byte, num),
		_size:     size, _num: num,
	}

	// for i := uint32(0); i < num; i++ {
	// 	b := bp.newBytes()
	// 	bp.Put(b)
	// }

	return
}

func (bp *BytePool) Get() (b []byte) {
	bp.Lock()

	if bp._readPos >= bp._writePos {
		// 扩容一倍
		num := bp._num
		bp._num = num * 2
		log.Printf("扩容到%d\n", bp._num)
		tmpNewBytesVec := make([][]byte, 0, bp._num)
		tmpNewBytesVec = append(tmpNewBytesVec, bp._bytesVec...)
		bp._writePos += num
		for i := num; i < bp._num; i++ {
			tmpNewBytesVec = append(tmpNewBytesVec, bp.newBytes())
		}
		bp._bytesVec = tmpNewBytesVec
	}

	idx := bp._readPos % bp._num
	bp._readPos++
	bp.Unlock()

	b = bp._bytesVec[idx]
	return
}

func (bp *BytePool) Put(b []byte) {
	bp.Lock()

	if bp._writePos < bp._readPos {
		panic("write pos small than read pos!")
	}

	idx := bp._writePos % bp._num
	bp._writePos++
	bp.Unlock()

	bp._bytesVec[idx] = b
}

func (bp *BytePool) Size() int { return int(bp._size) }

func (bp *BytePool) newBytes() (b []byte) {
	b = make([]byte, bp._size)
	return
}
