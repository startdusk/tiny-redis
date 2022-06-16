package aof

import (
	"os"
	"strconv"

	"github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/lib/utils"
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
