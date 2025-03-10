package netx

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
)

var (
	ErrPortTimeout = errors.New("waitPort.timeout")
)

func WaitForPorts(
	gtx context.Context,
	hostPort string,
	maxWait time.Duration) error {

	start := time.Now()
	_ = start
	open := false
	log.Info().Str("host", hostPort).Msg("waiting for port...")
	for start.Add(maxWait).After(time.Now()) {
		select {
		case <-gtx.Done():
			return gtx.Err()
		default:
		}
		if IsPortOpen(gtx, hostPort) {
			open = true
			break
		}

	}
	if !open {
		return errx.Errf(ErrPortTimeout, "timeout waiting for '%s'", hostPort)
	}

	return nil
}

func IsPortOpen(gtx context.Context, hostPort string) bool {
	conn, err := net.DialTimeout("tcp", hostPort, 1*time.Second)
	if err == nil && conn != nil {
		log.Info().Str("host", hostPort).Msg("port is now open")
		conn.Close()
		return true
	}
	return false
}
