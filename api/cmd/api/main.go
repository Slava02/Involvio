package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Slava02/Involvio/config"
	"github.com/Slava02/Involvio/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize the logger
	logger.SetupLogger(cfg)

	err = api.Run(ctx, cancel, cfg, slog.Default())
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to run application: %v", err))
		return
	}

}

func Run(ctx context.Context, cancelFunc context.CancelFunc, cfg *config.Config, logger *slog.Logger) error {
	// Run the application
	application := app.NewApp()
	go app.Run(application.Server, cfg)

	stopped := make(chan struct{})
	go func() {
		// Используем буферизированный канал, как рекомендовано внутри signal.Notify функции
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Блокируемся и ожидаем из канала quit - interrupt signal,
		// чтобы сделать gracefully shutdown с таймаутом в 10 сек
		<-quit

		// Завершаем работу горутин
		cancelFunc()

		// Получили SIGINT (0x2) или SIGTERM (0xf), выполняем graceful shutdown
		exitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := application.Server.ShutdownWithContext(exitCtx); err != nil {
			logger.Error("gracefully shutdown error")
		} else {
			logger.Warn("Server stopped")
		}

		close(stopped)
	}()

	<-stopped

	slog.Info("Server gracefully stopped, bye, bye!")

	return nil
}
