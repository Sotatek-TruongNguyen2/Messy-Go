package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	conn, err := pgxpool.New(context.Background(), "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable")

	if (err != nil) {
		log.Fatal("cannot connect to db: ", err)
	}

	testDB = conn
	testQueries = &Queries{db:conn}

	os.Exit(m.Run())
}