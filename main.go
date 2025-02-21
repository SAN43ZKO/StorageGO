package main

import (
	"Storage/internal/config"
	"Storage/internal/database"
	"Storage/internal/router"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func main() {

	logs := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	connStr := makeConnectionStr(cfg)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Connection Error: ", err)
	}
	fmt.Println("Connection successfuly")
	defer conn.Close(context.Background())

	if err := database.InitializeDB(conn); err != nil {
		log.Fatal("failed to initialize database error", err)

	}

	db := database.NewPostgresStorage(conn)
	producthandler := router.NewProducHandler(*db)
	registerRouter(producthandler)

	logs.Info("Starting server", "addres", cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port))
	if err := http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), nil); err != nil {
		logs.Error("Failed to start server", "err", err)
	}
}

func makeConnectionStr(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, strconv.Itoa(cfg.Postgres.Port), cfg.Postgres.DB)
}

func registerRouter(h *router.ProductHandler) {
	http.HandleFunc("GET /product/", h.GetProducts)
	http.HandleFunc("GET /product/{id}", h.GetProductsById)
	http.HandleFunc("DELETE /product/{id}", h.DeleteProducts)
	http.HandleFunc("POST /product", h.CreateProducts)
	http.HandleFunc("PUT /product/{id}", h.UpdateProducts)
}
