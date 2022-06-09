package resp

type Reply interface {
	Bytes() []byte
}
