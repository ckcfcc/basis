package pool

type IBytePool interface {
	Get() []byte
	Put([]byte)
	Size() int
}
