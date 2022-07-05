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
		<-sigs
		cancel()
		// code := 0
		// if sig == syscall.SIGTERM {
		// 	code = -1
		// }
		// log.Info().Str("signal", sig.String()).
		// 	Msg("received OS signal, exiting")
		// os.Exit(code)
	}()

	return gtx, cancel
}
