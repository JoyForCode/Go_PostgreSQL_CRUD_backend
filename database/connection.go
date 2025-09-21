package database

import (
	"context"
	"log"
	 "os"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func init() {
	err:=godotenv.Load(".env")
	if err!=nil {
		fmt.Print("Error loading the env")
	}
}

func Connect() *pgx.Conn {
	connection:=os.Getenv("CONNECTION_URL")
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return conn
}

func CreateUsersTable(conn *pgx.Conn) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := conn.Exec(context.Background(), query)
	return err
}

func TestConnection(conn *pgx.Conn) error {
	var result string
	err := conn.QueryRow(context.Background(), "SELECT 'Database is working!'").Scan(&result)
	return err
}
