package db_test

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"testing"

	. "github.com/mortum5/simple-bank/db/sqlc"

	_ "github.com/lib/pq"
	"github.com/mortum5/simple-bank/config"
	"golang.org/x/exp/slog"
)

const (
	dbDriver = "postgres"
)

var (
	testQueries  *Queries
	testDB       *sql.DB
	programLevel = new(slog.LevelVar)
)

func TestMain(m *testing.M) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: programLevel,
	}))
	slog.SetDefault(logger)
	programLevel.Set(-4)

	config, err := config.LoadConfig("../../")
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

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
