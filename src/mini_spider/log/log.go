package log

import (
	"path/filepath"

	"github.com/alecthomas/log4go"
)

var Logger log4go.Logger

var levelMap map[string]log4go.Level = map[string]log4go.Level{
	"DEBUG":    log4go.DEBUG,
	"TRACE":    log4go.TRACE,
	"INFO":     log4go.INFO,
	"WARNING":  log4go.WARNING,
	"ERROR":    log4go.ERROR,
	"CRITICAL": log4go.CRITICAL,
}

func Init(logFileNamePrefix, rootDirectory, levelStr string, hasStdout bool) {
	Logger = make(log4go.Logger)

	level := levelMap[levelStr]
	if level < levelMap["DEBUG"] {
		level = levelMap["INFO"]
	}

	writer := log4go.NewFileLogWriter(filepath.Join(rootDirectory, logFileNamePrefix+".log"), false)
	Logger.AddFilter("log", level, writer)

	writer = log4go.NewFileLogWriter(filepath.Join(rootDirectory, logFileNamePrefix+".wf.log"), false)
	Logger.AddFilter("wf_log", log4go.ERROR, writer)

	if hasStdout {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}
}
