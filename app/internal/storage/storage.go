/*
Структура бд:

-- users table
id
email
password_hash
created_at

-- habits table (описание привычек)
id
user_id
title
description
created_at

-- habit_progress (записи о выполнении привычек)
id
habit_id
copleted_date
created_at
unique(habit_id, completed_date)

-- sessions (stateful auth для возврата html-страниц)
id (uid)
user_id
created_at
expires_at

*/

package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/Artymka/habits-control/app/internal/config"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", getDsn(cfg))
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db:  db,
		log: logger,
	}
	err = storage.CreateTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func getDsn(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateTables() error {
	const op = "storage.create_tables"
	stmt := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS habits (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS habit_progress (
		id SERIAL PRIMARY KEY,
		habit_id INTEGER REFERENCES habits(id) ON DELETE CASCADE,
		completed_date DATE NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(habit_id, completed_date)
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		crea DATE NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		expires_at TIMESTAMP NOT NULL,
		UNIQUE(user_id)
	);
	`

	_, err := s.db.Exec(stmt)
	s.log.Debug("done", slog.String("op", op))
	return err
}

func (s *Storage) DropTables() error {
	const op = "storage.drop_tables"
	stmt := `
		DROP TABLE IF EXISTS habits_progress CASCADE;
		DROP TABLE IF EXISTS habits CASCADE;
		DROP TABLE IF EXISTS sessions CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
	`

	_, err := s.db.Exec(stmt)
	s.log.Debug("done", slog.String("op", op))
	return err
}
