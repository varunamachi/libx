package rt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Gtx() (context.Context, context.CancelFunc) {

	gtx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-sigs:
				cancel()
			case <-gtx.Done():
				return
			}
		}
	}()

	return gtx, cancel
}
