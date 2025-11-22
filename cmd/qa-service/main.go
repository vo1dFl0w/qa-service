package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vo1dFl0w/qa-service/internal/app/adapters/storage/postgres"
	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/config"
	"github.com/vo1dFl0w/qa-service/internal/logger"
	httptransport "github.com/vo1dFl0w/qa-service/internal/transport/http_transport"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "sturtup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// ===== LoadConfig =====
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	// ===== LoadLogger =====
	log := logger.LoadLogger(cfg.Env)

	// ===== Create databaseDSN by lodaded config =====
	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Sslmode,
	)

	// ===== Trying to connect to database by gorm with retries =====
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(pg.Open(databaseDSN))
		if err == nil {
			break
		}
		log.Info("waiting for database...")
    	time.Sleep(time.Second * 1)
	}
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	log.Info("connection to the database established")

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	defer sqlDB.Close()

	// ===== Create storage with usecases =====
	storage := postgres.New(db)
	questionUsecase := usecase.NewQuestionService(storage.Question())
	answerUsecase := usecase.NewAnswerService(storage.Answer(), storage.Question())

	// ===== Run the http server =====
	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httptransport.NewHandler(log, questionUsecase, answerUsecase),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		log.Info("server started", "host", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()
	
	// ===== Gracefull shutdown / server error detected =====
	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case s := <-sig:
		log.Info("initialization gracefull shutdown", "signal", s)
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		log.Info("server gracefully stopped")
		return nil
	}
}
