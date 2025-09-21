package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	err:=godotenv.Load(".env")
	if err!=nil {
		fmt.Print("Error loading the env")
	}
}

func generateConnString() (string, error) {
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_URL_OR_IP")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")
    
    if user == "" || password == "" || host == "" || port == "" || dbname == "" {
        return "", fmt.Errorf("required environment variables not set")
    }
    
    connectionURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        user, password, host, port, dbname)
    
    return connectionURL, nil
}

//Need to Deprecate use of *pgx.Conn for Concurrency handling
func Connect() *pgx.Conn {
	connection:=os.Getenv("CONNECTION_URL")
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return conn
}

func PoolConnect() *pgxpool.Pool {
	connectionURL, err:=generateConnString()
	if err!=nil{
		log.Fatalf("Database configuration error %s", err)
	}

	ctx, cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err:=pgxpool.ParseConfig(connectionURL)
	if err!=nil{
		log.Fatalf("Failed to parse DB config %s", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err:=conn.Exec(ctx, 
		`SET TIME ZONE 'UTC';
		SET search_path TO public;`)
		return err
	}

	pool,err:=pgxpool.NewWithConfig(ctx, config)
	if err!=nil {
		log.Fatalf("Failed to create pool:%s",err)
	}

	return pool
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
