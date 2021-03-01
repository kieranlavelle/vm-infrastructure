package api

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func connectToDatabase() *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func getAPIProxy(application string, conn *pgxpool.Pool) pgx.Row {
	query := `
		SELECT
			container_name, container_port
		FROM
			applications
		WHERE
			application=$1
	`
	return conn.QueryRow(context.Background(), query, application)
}
