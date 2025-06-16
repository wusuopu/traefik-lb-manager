package provider

import (
	"log"
	"os"
)

var _logger *log.Logger
var _logFile *os.File


func InitLogger(filename string) {
	if _logger != nil {
		// 已经初始化
		return
	}

	_logFile, _ = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if _logFile == nil {
		_logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		_logger = log.New(_logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func Logger() *log.Logger {
	return _logger
}

func CloseLogger() {
	if _logFile != nil {
		_logFile.Close()
		_logFile = nil
		_logger = nil
	}
}