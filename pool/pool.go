package pool

import (
	"github.com/ckcfcc/mathx"
)

var bytesMgr = &BytesMgr{}

func GetBytes(size int) []byte {
	return bytesMgr.Get(size)
}

type BytesMgr struct {
}

func (bm *BytesMgr) Get(size int) (b []byte) {
	p2s := Cale(size)

	return b
}
