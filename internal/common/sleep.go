package common

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
