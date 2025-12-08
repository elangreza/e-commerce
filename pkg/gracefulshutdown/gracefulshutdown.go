package gracefulshutdown

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Operation struct {
	Name         string
	ShutdownFunc func(ctx context.Context) error
}

func New(ctx context.Context, timeout time.Duration, ops ...Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		slog.Info("shutting down")

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		go func() {
			<-ctx.Done()
			slog.Info("force quit the app")
			wait <- struct{}{}
		}()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			go func(key int, op Operation) {
				defer wg.Done()

				slog.Info(op.Name, "shutdown", "started")

				if err := op.ShutdownFunc(ctx); err != nil {
					slog.Error(op.Name, "err", err.Error())
					return
				}

				slog.Info(op.Name, "shutdown", "finished")
			}(key, op)
		}

		wg.Wait()
	}()

	return wait
}
