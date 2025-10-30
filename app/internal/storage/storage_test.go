package storage

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Artymka/habits-control/app/internal/config"
	"github.com/Artymka/habits-control/app/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	ConfigPath = "./../../../config/local.yaml"
)

func TestDropCreateTables(t *testing.T) {
	t.Run("Create tables with no errors", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		cfg, err := config.New(ConfigPath)
		assert.Nil(t, err)

		s, err := New(cfg, logger)
		assert.Nil(t, err)

		err = s.DropTables()
		assert.Nil(t, err)

		err = s.CreateTables()
		assert.Nil(t, err)
	})
}

func TestOperations(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := config.New(ConfigPath)
	assert.Nil(t, err)
	s, err := New(cfg, logger)
	assert.Nil(t, err)

	var userID int64

	t.Run("Plain creation of user", func(t *testing.T) {
		var err error
		userID, err = s.CreateUser(models.UserCreate{
			Email:        "example@google.com",
			PasswordHash: "asdjf;alskdjf;alsdj",
		})

		t.Logf("created user id: %d\n", userID)
		assert.Nil(t, err)
	})

	t.Run("Plain user get", func(t *testing.T) {
		user, err := s.GetUser(userID)

		assert.Nil(t, err)
		assert.Equal(t, "example@google.com", user.Email)
	})
}
