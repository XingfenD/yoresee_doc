package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdownSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}
