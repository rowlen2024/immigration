package handler

import (
	"os"
	"testing"

	"mygo-immigration/backend/internal/logging"
)

func TestMain(m *testing.M) {
	logging.Init("info")
	os.Exit(m.Run())
}
