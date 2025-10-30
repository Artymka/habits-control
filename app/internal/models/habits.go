package models

import "time"

type HabitCreate struct {
	UserID      int64
	Title       string
	Description string
}

type HabitUpdate struct {
	ID          int64
	Title       string
	Description string
}

type Habit struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	CreatedAt   time.Time
}

type HabitResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type HabitThinResponse struct {
	ID int64 `json:"id"`
}

type HabitRequest struct {
	Title       string `json:"title" validate:"min=3,max=255"`
	Description string `json:"description"`
}
