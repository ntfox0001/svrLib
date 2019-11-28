package log

// 修改来自"github.com/inconshreveable/log15"
// 主要修改 不再过滤换行符

import (
	"os"
)

func Initial(logFormat string, level Lvl, logFile string) error {

	// 默认文件格式是fmt
	var logf Format
	switch logFormat {
	case "json":
		logf = JsonFormatEx(false, true)
	case "std":
		logf = TerminalFormat()
	default: // "fmt":
		logf = LogfmtFormat()

	}

	// os标准输出总是terminal
	var comboHandler Handler
	stdHandler := StreamHandler(os.Stdout, logf)
	if logFile != "" {
		fileHandler := Must.FileHandler(logFile, logf)
		comboHandler = MultiHandler(stdHandler, fileHandler)
	} else {
		comboHandler = stdHandler
	}

	Root().SetHandler(CallerFileHandler(LvlFilterHandler(level, comboHandler)))

	return nil
}
