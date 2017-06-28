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

func ConfigureLogger(file, terminal bool, fileLevel, terminalLevel logging.Level, ) {

	var backendList []logging.Backend

	if terminal {
		// Set a custom formatting for the logs sent to the terminal
		terminalFormat := logging.MustStringFormatter("%{color}[%{shortfunc}]%{color:reset} " +
			"%{time:15:04:05.000} %{color}%{level:.4s}%{color:reset} -> %{message}")
		terminalBackend := logging.NewLogBackend(os.Stdout, "", 0)

		// Only messages of the specified minimum level or higher should be sent to the terminal backend
		terminalBackendLeveled := logging.AddModuleLevel(terminalBackend)
		terminalBackendLeveled.SetLevel(terminalLevel, "")
		terminalBackendFormatter := logging.NewBackendFormatter(terminalBackendLeveled, terminalFormat)

		backendList = append(backendList, terminalBackendFormatter)
	}

	if file {
		// Set a custom formatting for the logs sent to the terminal
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
		fileBackendLeveled := logging.AddModuleLevel(fileBackend)
		fileBackendLeveled.SetLevel(fileLevel, "")
		fileBackendFormatter := logging.NewBackendFormatter(fileBackendLeveled, fileFormat)

		backendList = append(backendList, fileBackendFormatter)
	}

	// Set the enabled backends for the logger
	logging.SetBackend(backendList...)
}
