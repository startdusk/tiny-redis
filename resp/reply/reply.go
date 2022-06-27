package reply

import (
	"bytes"
	"strconv"

	"github.com/startdusk/tiny-redis/api/resp"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type BulkReply struct {
	arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{arg: arg}
}

func (r *BulkReply) Bytes() []byte {
	if len(r.arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte(`$` + strconv.Itoa(len(r.arg)) + CRLF + string(r.arg) + CRLF)
}

type MultiBulkReply struct {
	args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{args: args}
}

func (r *MultiBulkReply) Args() [][]byte {
	return r.args
}

func (r *MultiBulkReply) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(len(r.args)) + CRLF)
	for _, arg := range r.args {
		if arg == nil {
			buf.WriteString(string(nullBulkReply) + CRLF)
		} else {
			buf.WriteString(`$` + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

type StatusReply struct {
	status string
}

func NewStatusReply(status string) *StatusReply {
	return &StatusReply{status: status}
}

func (r *StatusReply) Bytes() []byte {
	return []byte("+" + r.status + CRLF)
}

type NumberReply struct {
	Code int64
}

func NewNumberReply(code int64) *NumberReply {
	return &NumberReply{Code: code}
}

func (r *NumberReply) Bytes() []byte {
	return []byte(":" + strconv.FormatInt(int64(r.Code), 10) + CRLF)
}

type ErrorReply interface {
	Error() string
	Bytes() []byte
}

type StandardErrReply struct {
	status string
}

func NewStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{status: status}
}

func (r *StandardErrReply) Bytes() []byte {
	return []byte("-" + r.status + CRLF)
}

func (r *StandardErrReply) Error() string {
	return r.status
}

func IsErrRely(reply resp.Reply) bool {
	if len(reply.Bytes()) == 0 {
		return false
	}
	return reply.Bytes()[0] == '-'
}
