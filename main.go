package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ymgyt/ops/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watchSignal(cancel)
	cmd.New().Execute(ctx)
}

func watchSignal(cancel func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	cancel()
}
