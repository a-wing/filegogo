package server

import (
	"context"
	"time"
)

func run(ctx context.Context, fn func()) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
			case <- ticker.C:
				fn()
			case <- ctx.Done():
				return
		}
	}
}
