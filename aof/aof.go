package aof

import (
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/lib/utils"
	"github.com/startdusk/tiny-redis/resp/conn"
	"github.com/startdusk/tiny-redis/resp/parser"
	"github.com/startdusk/tiny-redis/resp/reply"
)

type CmdLine = [][]byte

const aofBufferSize = 1 << 16 // 65535

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type Handler struct {
	db         db.Database
	payloadCh  chan *payload
	file       *os.File
	filename   string
	currentDB  int
	appendOnly bool
}

func NewAOFHandler(db db.Database, aofFilename string, appendOnly bool) (*Handler, error) {
	h := Handler{db: db, filename: aofFilename, appendOnly: appendOnly}
	h.Load()
	file, err := os.OpenFile(h.filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	h.file = file
	if h.appendOnly {
		h.payloadCh = make(chan *payload, aofBufferSize)
		go h.handle()
	}
	return &h, nil
}

func (h *Handler) Add(dbIndex int, cmdLine CmdLine) {
	if h.appendOnly && h.payloadCh != nil {
		h.payloadCh <- &payload{
			dbIndex: dbIndex,
			cmdLine: cmdLine,
		}
	}
}

func (h *Handler) Load() {
	file, err := os.Open(h.filename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()
	stream := parser.ParseStream(file)
	fakeConn := &conn.Conn{}
	for p := range stream {
		if p.Err != nil {
			if errors.Is(p.Err, io.EOF) {
				break
			}
			logger.Error(err)
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi bulk")
			continue
		}

		rly := h.db.Exec(fakeConn, r.Args())
		if reply.IsErrRely(rly) {
			logger.Error(rly)
		}
	}
}

func (h *Handler) handle() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	for p := range h.payloadCh {
		if p.dbIndex != h.currentDB {
			r := reply.NewMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex)))
			_, err := h.file.Write(r.Bytes())
			if err != nil {
				logger.Error(err)
				continue
			}
			h.currentDB = p.dbIndex
		}

		r := reply.NewMultiBulkReply(p.cmdLine)
		_, err := h.file.Write(r.Bytes())
		if err != nil {
			logger.Error(err)
		}
	}
}
