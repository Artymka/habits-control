package storage

import (
	"log/slog"

	"github.com/Artymka/habits-control/app/internal/models"
)

func (s *Storage) CreateUser(user models.UserCreate) (int64, error) {
	const op = "storage.create_user"
	stmt, err := s.db.Prepare(`
		INSERT INTO users
		(email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`)

	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRow(user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, err
	}

	s.log.Debug("done", slog.String("op", op))
	return id, nil
}

func (s *Storage) GetUser(id int64) (models.User, error) {
	const op = "storage.get_user"
	query, err := s.db.Prepare(`
		SELECT id, email, password_hash, created_at FROM users WHERE id = $1
	`)
	if err != nil {
		return models.User{}, err
	}

	res := models.User{}
	err = query.QueryRow(id).Scan(&res.ID, &res.Email, &res.PasswordHash, &res.CreatedAt)

	s.log.Debug("done", slog.String("op", op))
	return res, err
}
