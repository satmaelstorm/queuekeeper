package main

import (
	"fmt"
	"log"
	"os"
)

const (
	QK_LOG_LEVEL_TRACE    = 10
	QK_LOG_LEVEL_DEBUG    = 20
	QK_LOG_LEVEL_INFO     = 30
	QK_LOG_LEVEL_WARNING  = 40
	QK_LOG_LEVEL_ERROR    = 50
	QK_LOG_LEVEL_CRITICAL = 60
)

const (
	QK_DEFAULT_LOG_LEVEL  = QK_LOG_LEVEL_WARNING
	QK_DEFAULT_LOG_ENGINE = "std"
)

const (
	QK_STR_LOG_LEVEL_TRACE    = "TRACE"
	QK_STR_LOG_LEVEL_DEBUG    = "DEBUG"
	QK_STR_LOG_LEVEL_INFO     = "INFO"
	QK_STR_LOG_LEVEL_WARNING  = "WARNING"
	QK_STR_LOG_LEVEL_ERROR    = "ERROR"
	QK_STR_LOG_LEVEL_CRITICAL = "CRITICAL"
)

type queueKeeperLogger struct {
	logger *log.Logger
}

func parseLogLevel(level string) int {
	l := QK_DEFAULT_LOG_LEVEL
	switch level {
	case QK_STR_LOG_LEVEL_CRITICAL:
		l = QK_LOG_LEVEL_CRITICAL
	case QK_STR_LOG_LEVEL_DEBUG:
		l = QK_LOG_LEVEL_DEBUG
	case QK_STR_LOG_LEVEL_ERROR:
		l = QK_LOG_LEVEL_ERROR
	case QK_STR_LOG_LEVEL_INFO:
		l = QK_LOG_LEVEL_INFO
	case QK_STR_LOG_LEVEL_TRACE:
		l = QK_LOG_LEVEL_TRACE
	case QK_STR_LOG_LEVEL_WARNING:
		l = QK_LOG_LEVEL_WARNING
	}
	return l
}

func parseLogLevelToString(level int) string {
	l := "UNKNOWN"
	switch level {
	case QK_LOG_LEVEL_CRITICAL:
		l = QK_STR_LOG_LEVEL_CRITICAL
	case QK_LOG_LEVEL_DEBUG:
		l = QK_STR_LOG_LEVEL_DEBUG
	case QK_LOG_LEVEL_ERROR:
		l = QK_STR_LOG_LEVEL_ERROR
	case QK_LOG_LEVEL_INFO:
		l = QK_STR_LOG_LEVEL_INFO
	case QK_LOG_LEVEL_TRACE:
		l = QK_STR_LOG_LEVEL_TRACE
	case QK_LOG_LEVEL_WARNING:
		l = QK_STR_LOG_LEVEL_WARNING
	}
	return l
}

func initLogger(conf logConfiguration) queueKeeperLogger {
	var logger = queueKeeperLogger{}
	if QK_DEFAULT_LOG_ENGINE == conf.engine || "" == conf.parsedEngine.Scheme {
		logger.logger = log.New(os.Stderr, "QueueKeeper", log.LstdFlags|log.Lshortfile)
	} else {
		panic("Now implemeted only 'std' engine")
	}
	return logger
}

func (qkl queueKeeperLogger) log(level int, message string) {
	message = fmt.Sprintf("[%s] %s", parseLogLevelToString(level), message)
	qkl.logger.Println(message)
}
