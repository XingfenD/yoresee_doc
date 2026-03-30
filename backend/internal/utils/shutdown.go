package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

type ShutdownHook struct {
	Name string
	Fn   func() error
}

func ShutdownContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func WaitForShutdownSignal() {
	ctx, stop := ShutdownContext()
	defer stop()
	<-ctx.Done()
}

func RunShutdownHooks(hooks ...ShutdownHook) {
	for _, hook := range hooks {
		if hook.Fn == nil {
			continue
		}
		if err := hook.Fn(); err != nil {
			logrus.Errorf("Close %s failed: %v", hook.Name, err)
		}
	}
}
