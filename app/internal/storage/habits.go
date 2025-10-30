package storage

import (
	"log/slog"

	"github.com/Artymka/habits-control/app/internal/models"
)

func (s *Storage) GetHabitsOfUser(userID int64) ([]models.HabitResponse, error) {
	stmt, err := s.db.Prepare(`
		SELECT id, user_id, title, description, created_at
		FROM habits
		WHERE user_id = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	res := make([]models.HabitResponse, 0)
	for rows.Next() {
		res = append(res, models.HabitResponse{})
		i := len(res) - 1
		if err = rows.Scan(
			&res[i].ID,
			&res[i].Title,
			&res[i].Description,
			&res[i].CreatedAt); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (s *Storage) CreateHabit(habit models.HabitCreate) (int64, error) {
	const op = "storage.create_habit"
	stmt, err := s.db.Prepare(`
		INSERT INTO habits
		(user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`)
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRow(habit.UserID, habit.Title, habit.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	s.log.Debug("done", slog.String("op", op))
	return id, nil
}

func (s *Storage) UpdateHabit(habit models.HabitUpdate) error {
	const op = "storage.update_habit"
	stmt, err := s.db.Prepare(`
		UPDATE habits SET
		title = $1 description = $2
		WHERE id = $3
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(habit.Title, habit.Description, habit.ID)
	if err != nil {
		return err
	}

	s.log.Debug("done", slog.String("op", op))
	return nil
}

func (s *Storage) DeleteHabit(id int64) error {
	const op = "delete_habit"
	stmt, err := s.db.Prepare(`
		DELETE FROM habits
		WHERE id = $1
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	s.log.Debug("done", slog.String("op", op))
	return nil
}
