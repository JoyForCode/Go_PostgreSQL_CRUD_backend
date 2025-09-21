package database

import (
	"context"
	"postgre_advanced/models"

	"github.com/jackc/pgx/v5"
)

func CreateUser(conn *pgx.Conn, req models.CreateUserRequest) error {
	query := `INSERT INTO users (name, email) values ($1, $2)`
	_, err := conn.Exec(context.Background(), query, req.Name, req.Email)
	return err
}

func GetAllUsers(conn *pgx.Conn) ([]models.User, error) {
	query := `SELECT id,name,email,created_at FROM users`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return users, nil
}

func GetUserByID(conn *pgx.Conn, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id,name,email,created_at FROM users WHERE id=$1`
	err := conn.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(conn *pgx.Conn, id int, req models.UpdateUserRequest) error {
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	result, err := conn.Exec(context.Background(), query, req.Name, req.Email, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func DeleteUser(conn *pgx.Conn, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	result, err := conn.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
