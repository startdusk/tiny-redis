package reply

type UnknownErrReply struct{}

func NewUnknowErrReply() *UnknownErrReply {
	return newUnknownErrReply
}

var unknownErrReply = []byte("-Err unknown\r\n")
var newUnknownErrReply = new(UnknownErrReply)

func (r UnknownErrReply) Bytes() []byte {
	return unknownErrReply
}

func (r UnknownErrReply) Error() string {
	return "Err nuknown"
}

type ArgNumErrReply struct {
	cmd string
}

func NewArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{cmd: cmd}
}

func (r *ArgNumErrReply) Bytes() []byte {
	return []byte(`-Err wrong number of arguments for "` + r.cmd + `" command\r\n`)
}

func (r *ArgNumErrReply) Error() string {
	return `-Err wrong number of arguments for "` + r.cmd + `" command`
}

type SyntaxErrReply struct{}

func NewSyntaxErrReply() *SyntaxErrReply {
	return newSyntaxErrReply
}

var syntaxErrReply = []byte("-Err syntax error\r\n")
var newSyntaxErrReply = new(SyntaxErrReply)

func (r SyntaxErrReply) Bytes() []byte {
	return syntaxErrReply
}

func (r SyntaxErrReply) Error() string {
	return "syntax error"
}

type WrongTypeErrReply struct{}

func NewWrongTypeErrReply() *WrongTypeErrReply {
	return newWrongTypeErrReply
}

var wrongTypeErrReply = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
var newWrongTypeErrReply = new(WrongTypeErrReply)

func (r WrongTypeErrReply) Bytes() []byte {
	return wrongTypeErrReply
}

func (r WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

type ProtocolErrReply struct {
	msg string
}

func NewProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{msg: msg}
}

func (r ProtocolErrReply) Bytes() []byte {
	return []byte(`-Err protocol error "` + r.msg + `"\r\n`)
}

func (r ProtocolErrReply) Error() string {
	return `Err protocol error "` + r.msg + `"`
}
