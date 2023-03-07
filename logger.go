package firespace

import "go.uber.org/zap"

var logger = zap.NewNop()
var sugar = logger.Sugar()

func GetLogger() *zap.SugaredLogger {
	return sugar
}

func SetLogger(newSugar *zap.SugaredLogger) {
	sugar = newSugar
}
