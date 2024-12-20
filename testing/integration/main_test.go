package integration

import (
	"os"
	"testing"

	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/logger"
)

func TestMain(m *testing.M) {
	logger.SetLogLevels(logger.LevelDebug)
	logger.InitLogStdout(logger.LevelDebug)

	os.Exit(m.Run())
}
