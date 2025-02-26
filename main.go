package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"Storage/internal/config"
	"Storage/internal/handler"
	"Storage/internal/storage"
	"Storage/pkg/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("app", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("run (1): %w", err)
	}

	db, err := postgres.NewPostgres(ctx, cfg)
	if err != nil {
		return fmt.Errorf("run (2): %w", err)
	}
	defer db.Close(ctx)

	storage, err := storage.NewStorage(ctx, logger, db)
	if err != nil {
		return fmt.Errorf("run (3): %w", err)
	}

	if err := storage.GetAllProduct(ctx, db); err != nil {
		return fmt.Errorf("run (4): %w", err)
	}

	go storage.CacheUpdater(ctx)

	productHandler := handler.NewProducHandler(logger, storage)

	router := http.NewServeMux()

	router.HandleFunc("GET /product/", productHandler.GetProducts)
	router.HandleFunc("GET /product/{id}", productHandler.GetProductsById)
	router.HandleFunc("DELETE /product/{id}", productHandler.DeleteProducts)
	router.HandleFunc("POST /product", productHandler.CreateProducts)
	router.HandleFunc("PUT /product/{id}", productHandler.UpdateProducts)

	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router,
	}

	logger.Info("starting http server", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run (5): %w", err)
	}

	return nil
}
