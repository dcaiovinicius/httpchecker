package logger

import (
	"io"
	"log"
	"os"
)

var (
	infoLog   = log.New(io.Discard, "INFO: ", log.LstdFlags)
	errorLog  = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
	noticeLog = log.New(os.Stdout, "", log.LstdFlags)
)

func EnableDebug() {
	infoLog.SetOutput(os.Stdout)
}

func DisableDebug() {
	infoLog.SetOutput(io.Discard)
}

func Info(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Error(format string, v ...any) {
	errorLog.Printf(format, v...)
}

func Notice(format string, v ...any) {
	noticeLog.Printf(format, v...)
}
