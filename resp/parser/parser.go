package parser

import (
	"io"

	"github.com/startdusk/tiny-redis/api/resp"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && s.expectedArgsCount == len(s.args)
}

func ParseStream(r io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(r, ch)
	return ch
}

func parse0(r io.Reader, ch chan<- *Payload) {
}
