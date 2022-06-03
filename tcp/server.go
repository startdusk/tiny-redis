package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/startdusk/tiny-redis/api/tcp"
	"github.com/startdusk/tiny-redis/lib/logger"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeCh := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			close(closeCh)
		}
	}()
	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	return ListenAndServe(lis, handler, closeCh)
}

func ListenAndServe(lis net.Listener, handler tcp.Handler, closeCh <-chan struct{}) error {
	go func() {
		<-closeCh
		logger.Info("shutting down")
		lis.Close()
		handler.Close()
	}()

	defer func() {
		lis.Close()
		handler.Close()
	}()
	
	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.Error("accepted error", err)
			break
		}
		logger.Info("accepted link")
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				conn.Close()
			}()
			handler.Handle(ctx, conn)
		}()
	}

	wg.Wait()
	return nil
}
