package logger

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
)

var (
	InfoLogger    *log.Logger
	WarnLogger    *log.Logger
	SuccessLogger *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, constants.BlueColorOutput+"[INFO] ", log.Ldate|log.Ltime)
	WarnLogger = log.New(os.Stdout, constants.YellowColorOutput+"[WARN] ", log.Ldate|log.Ltime)
	SuccessLogger = log.New(os.Stderr, constants.GreenColorOuput+"[SUCCESS] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, constants.RedColorOuput+"[ERROR] ", log.Ldate|log.Ltime)
}

func logMessage(logger *log.Logger, format string, args ...interface{}) {
	functionInvocationLevelsUp := 3
	_, file, line, _ := runtime.Caller(functionInvocationLevelsUp)
	filename := path.Base(file)
	logger.Printf("[%s:%d] "+format+constants.DefaultColorOutput, append([]interface{}{filename, line}, args...)...)
}

func Info(message interface{}) {
	Log(InfoLogger, message)
}

func Warn(message interface{}) {
	Log(WarnLogger, message)
}

func Success(message interface{}) {
	Log(SuccessLogger, message)
}

func Error(message interface{}) {
	Log(ErrorLogger, message)
}

func Log(logger *log.Logger, message interface{}) {
	switch msg := message.(type) {
	case string:
		logMessage(logger, msg)
	default:
		jsonString, err := json.Marshal(msg)
		if err != nil {
			logMessage(logger, "%v", msg)
		}

		logMessage(logger, "%s", jsonString)
	}
}
