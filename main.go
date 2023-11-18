package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/mortum5/simple-bank/api"
	"github.com/mortum5/simple-bank/config"
	db "github.com/mortum5/simple-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	hostAndPort := net.JoinHostPort(config.DBHost, strconv.Itoa(config.DBPort))
	dbSource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPass,
		hostAndPort,
		config.DBName,
	)

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
