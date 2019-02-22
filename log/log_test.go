package log_test

import (
	"github.com/ntfox0001/svrLib/log"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	h1 := log.CallerFileHandler(log.StreamHandler(os.Stdout, log.TerminalFormat()))

	log.Root().SetHandler(h1)

	log.Debug("Debug", "key", "value")
	log.Info("Info", "key", "value")
	log.Warn("Warn", "key", "value")
	log.Error("Error", "key", "value")
	log.Crit("Crit", "key", "value")
}
