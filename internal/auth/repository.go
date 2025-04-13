package auth

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) createUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error generating password hash")
	}

	query := `INSERT INTO users(name, email, login, password) 
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`

	return r.db.QueryRow(query, user.Name, user.Email, user.Login, hashedPassword).
		Scan(&user.ID)
}

func (r *AuthRepository) GetUserUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Login, &user.Password)

	if err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}
