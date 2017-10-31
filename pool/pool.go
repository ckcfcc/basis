package pool

import (
	"github.com/ckcfcc/basis/mathx"
)

const MinPO2 = 5  // 32 Byte
const MaxPo2 = 26 // 64 MByte
const DPoolNum = 64

var dPoolNum int = 64
var bytesMgr = &BytesMgr{}

func SetPoolNum(v int) {
	if v > DPoolNum {
		dPoolNum = v
	}
}

func GetBytes(size int) []byte {
	return bytesMgr.Get(size)
}

type BytesMgr struct {
	_bytesMap [MaxPo2 - MinPO2]*BytePool
}

func (bm *BytesMgr) Get(size int) (b []byte) {
	p2s := mathx.MinPowerOf2(size)

	return b
}
