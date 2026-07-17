package repository

import (
	"database/sql"
	"go-ecommerce-api/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result, err := r.db.Exec("INSERT INTO users(name, email, password_hash) VALUES (?, ?, ?)", user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	row := r.db.QueryRow("SELECT id, name, email, password_hash from users where email = ?", email)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	row := r.db.QueryRow("SELECT id, name, email from users where id = ?", id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
