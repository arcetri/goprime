package rieseltest

import (
	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"fmt"
	"time"
)

// Initialize logger to be used in this package
var log = logging.MustGetLogger("rieseltest")

// ConfigureLogger allows to enable file and terminal outputs for the log messages
// The available logging levels are:
//
// 		DEBUG < INFO < NOTICE < WARNING < ERROR < CRITICAL
//
// It is not recommended to set the level of any backend to DEBUG unless for very specific
// debugging reasons. If you do so, expect a lot of logs and a considerable slowing down in performance.
func ConfigureLogger(file bool, fileLevel logging.Level, terminal bool, terminalLevel logging.Level) {
	var backendList []logging.Backend

	if terminal {
		// Set a custom formatting for the logs sent to the terminal
		terminalFormat := logging.MustStringFormatter("%{color}[%{shortfunc}]%{color:reset} " +
			"%{time:15:04:05.000} %{color}%{level:.4s}%{color:reset} -> %{message}")
		terminalBackend := logging.NewLogBackend(os.Stdout, "", 0)
		terminalBackendFormatted := logging.NewBackendFormatter(terminalBackend, terminalFormat)

		// Only messages of the specified minimum level or higher should be sent to the terminal backend
		terminalBackendLeveled := logging.AddModuleLevel(terminalBackendFormatted)
		terminalBackendLeveled.SetLevel(terminalLevel, "")

		backendList = append(backendList, terminalBackendLeveled)
	}

	if file {
		// Set a custom formatting for the logs sent to files
		fileFormat := logging.MustStringFormatter("[%{pid} - %{shortfunc}] %{time:15:04:05.000} " +
			"%{level:.4s} -> %{message}")

		// Let the file logs be saved on a rolling basis
		fileBackend := logging.NewLogBackend(&lumberjack.Logger{
			Filename:   fmt.Sprintf(".logs/logFile.%d%.2d%.2d.log",
				time.Now().Year(), time.Now().Month(), time.Now().Day()),
			MaxSize:    100, // megabytes
			MaxBackups: 2,
			MaxAge:     3, //days
		}, "", 0)

		// Only messages of the specified minimum level or higher should be sent to the file backend
		fileBackendFormatted := logging.NewBackendFormatter(fileBackend, fileFormat)
		fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatted)
		fileBackendLeveled.SetLevel(fileLevel, "")

		backendList = append(backendList, fileBackendLeveled)
	}

	// Set the enabled backends for the logger
	logging.SetBackend(backendList...)
}
