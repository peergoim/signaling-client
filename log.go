package signaling_client

import (
	"fmt"
	"sync/atomic"
)

var (
	// 0: debug
	debugLevel int32 = 0
	// 1: info
	infoLevel int32 = 1
	// 2: warn
	warnLevel int32 = 2
	// 3: error
	ErrorLevel int32 = 3
)

var (
	logLevel = atomic.AddInt32(&debugLevel, 0)
)

func setLogLevel(level string) {
	switch level {
	case "debug":
		atomic.StoreInt32(&logLevel, debugLevel)
	case "info":
		atomic.StoreInt32(&logLevel, infoLevel)
	case "warn":
		atomic.StoreInt32(&logLevel, warnLevel)
	case "error":
		atomic.StoreInt32(&logLevel, ErrorLevel)
	default:
		atomic.StoreInt32(&logLevel, debugLevel)
	}
}

const (
	// 不同的日志级别对应不同的颜色
	// 0: debug
	debugColor = "\033[0;34mDEBU\033[0m "
	// 1: info
	infoColor = "\033[0;32mINFO\033[0m "
	// 2: warn
	warnColor = "\033[0;33mWARN\033[0m "
	// 3: error
	errorColor = "\033[0;31mERRO\033[0m "
)

func logf(color, format string, args ...interface{}) {
	format = color + format + "\n"
	fmt.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	if logLevel <= debugLevel {
		logf(debugColor, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel <= infoLevel {
		logf(infoColor, format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logLevel <= warnLevel {
		logf(warnColor, format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel <= ErrorLevel {
		logf(errorColor, format, args...)
	}
}
