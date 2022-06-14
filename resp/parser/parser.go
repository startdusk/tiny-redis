package parser

import (
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/startdusk/tiny-redis/api/resp"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/resp/reply"
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
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	bufReader := bufio.NewReader(r)
	var state readState
	var msg []byte
	var ioErr bool
	var err error
	for {
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			ch <- &Payload{
				Err: err,
			}
			if ioErr {
				close(ch)
				return
			}
			state = readState{}
			continue
		}

		if !state.readingMultiLine {
			if msg[0] == '*' { // starts with *
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.NewEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // starts with $
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}

				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: reply.NewEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else {
				reply, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: reply,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: err,
				}
				state = readState{}
				continue
			}

			if state.finished() {
				var ry resp.Reply
				if state.msgType == '*' {
					ry = reply.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					ry = reply.NewBulkReply(state.args[0])
				} else {
					err = errors.New("protocol error: " + string(msg))
				}
				ch <- &Payload{
					Data: ry,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// example: *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 { // \r\n
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		msg = make([]byte, state.bulkLen+2) // plus len("\r\n")
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}

		// ends with "\r\n"
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		// read end
		state.bulkLen = 0
	}
	return msg, false, nil
}

func parseMultiBulkHeader(msg []byte, state *readState) error {
	if len(msg) < 4 {
		return errors.New("protocol error: " + string(msg))
	}
	// extract number
	// *3\r\n
	// $3\r\n
	expectedLine, err := strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.expectedArgsCount = int(expectedLine)
		state.readingMultiLine = true
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// example: $4\r\nPING\r\n
func parseBulkHeader(msg []byte, state *readState) (err error) {
	// extract number
	// example: $4\r\nPING\r\n
	if len(msg) < 4 {
		return errors.New("protocol error: " + string(msg))
	}
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// +OK\r\n -err\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	msgStr := string(msg)
	s := strings.TrimSuffix(msgStr, "\r\n")
	if len(s) <= 1 {
		return nil, errors.New("protocol error: " + msgStr)
	}
	switch s[0] {
	case '+':
		return reply.NewStatusReply(s[1:]), nil
	case '-':
		return reply.NewStandardErrReply(s[1:]), nil
	case ':':
		val, err := strconv.ParseInt(s[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + msgStr)
		}
		return reply.NewNumberReply(int(val)), nil
	default:
		return nil, errors.New("protocol error: " + msgStr)
	}
}

// example: *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
func readBody(msg []byte, state *readState) (err error) {
	if len(msg) < 3 {
		return errors.New("protocol error: " + string(msg))
	}
	line := msg[0 : len(msg)-2]

	// starts with $
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}

		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		// starts with text
		state.args = append(state.args, line)
	}
	return nil
}
