package initialize

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() {
	// build a development logger that record Debug And Above Logs to standard errors
	logger, err := zap.NewDevelopment()
	if err != nil {
		// Print + Exit
		log.Fatal("Logger Initial Failed", err.Error())
	}
	// set as the global logger, and can get the logger by "zap.L()" to record in anywhere
	zap.ReplaceGlobals(logger)
}
