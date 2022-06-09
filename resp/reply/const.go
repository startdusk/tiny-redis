package reply

type PongReply struct {
}

func NewPongReply() *PongReply {
	return newPongReply
}

var pongReply = []byte("+PONG\r\n")
var newPongReply = new(PongReply)

func (r PongReply) Bytes() []byte {
	return pongReply
}

type OKReply struct {
}

func NewOKReply() *OKReply {
	return newOKReply
}

var okReply = []byte("+OK\r\n")
var newOKReply = new(OKReply)

func (r OKReply) Bytes() []byte {
	return okReply
}

// empty string reply
type NullBulkReply struct{}

func NewNullBulkReply() *NullBulkReply {
	return newNullBulkReply
}

var nullBulkReply = []byte("$-1\r\n")
var newNullBulkReply = new(NullBulkReply)

func (r NullBulkReply) Bytes() []byte {
	return nullBulkReply
}

// empty array reply
type EmptyMultiBulkReply struct{}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return newEmptyMultiBulkReply
}

var emptyMultiBulkReply = []byte("*0\r\n")
var newEmptyMultiBulkReply = new(EmptyMultiBulkReply)

func (r EmptyMultiBulkReply) Bytes() []byte {
	return emptyMultiBulkReply
}

type NoReply struct{}

func NewNoReply() *NoReply {
	return newNoReply
}

var noReply = []byte("")
var newNoReply = new(NoReply)

func (r NoReply) Bytes() []byte {
	return noReply
}
