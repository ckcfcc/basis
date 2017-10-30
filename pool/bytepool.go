package pool

type BytePool struct {
	_bytes      chan []byte
	_size, _num int
}

func NewBytePool(size, num int) (bp *BytePool) {
	bp = &BytePool{
		_bytes: make(chan []byte, num),
		_size:  size, _num: num,
	}
	return
}

func (bp *BytePool) Get() (b []byte) {
	select {
	case b = <-bp._bytes:
	default:
		b = bp.newBytes()
	}
	return
}

func (bp *BytePool) Put(b []byte) {
	select {
	case bp._bytes <- b:
	default:
		// free b
	}
}

func (bp *BytePool) Size() int { return len(bp._bytes) }

func (bp *BytePool) newBytes() (b []byte) {
	b = make([]byte, bp._size)
	return
}
