package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func (a *App) startServer(ctx context.Context) {
	slog.Info("http server started", "addr", a.httpServer.Addr)

	if err := a.httpServer.ListenAndServe(); err != nil {
		a.errC <- fmt.Errorf("server: %w", err)
	}
}

func (a *App) initShutdown(ctx context.Context) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT)
	<-signalCh

	a.cancelCtx()

	a.errC <- a.httpServer.Shutdown(ctx)
}
