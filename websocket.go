package main

import (
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

var (
	errConnClosed = errors.New("conn closed")
)

func dialWebsocketChan(ctx context.Context, url string) chan []byte {
	fallbackMaxSec := 64
	initialSec := 1
	ch := make(chan []byte)

	go func() {
	Out:
		for {
			for initialSec <= fallbackMaxSec {
				time.Sleep(time.Duration(initialSec) * time.Second)
				cctx, cancel := context.WithCancel(ctx)
				if err := dialWebsocketToChan(cctx, url, ch); err == cctx.Err() {
					cancel()
					break Out
				} else if err == errConnClosed {
					initialSec = 1
				}

				cancel()

				if initialSec < fallbackMaxSec {
					initialSec *= 2
				}
			}
		}
	}()

	return ch
}

func dialWebsocketToChan(ctx context.Context, url string, ch chan []byte) error {
	dialer := &websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}

	// ping pong
	go func() {
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return ctx.Err()
		default:
			_, buf, err := conn.ReadMessage()
			if err != nil {
				break Loop
			}
			ch <- buf
		}
	}

	return errConnClosed
}
