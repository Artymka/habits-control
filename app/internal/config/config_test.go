package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Plain load", func(t *testing.T) {
		wd, _ := os.Getwd()
		t.Logf("working directory: %s\n", wd)
		_, err := New("./../../../config/local.yaml")
		assert.Nil(t, err)
		// t.Log(cfg)
	})
}
