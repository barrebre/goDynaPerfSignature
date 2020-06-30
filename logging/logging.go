package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

var logLevel = "ERROR"

// SetLogLevel allows the app to set the appropriate log level based on config
func SetLogLevel(level string) {
	logLevel = level
}

// LogSystem will log with level ERROR
func LogSystem(info datatypes.Logging) {
	writeLog(info, "SYSTEM")
}

// LogError will log with level ERROR
func LogError(info datatypes.Logging) {
	writeLog(info, "ERROR")
}

// LogDebug will log with level Info if the logging level is set as such
func LogDebug(info datatypes.Logging) {
	level := getLoggingLevel()
	if level == "DEBUG" {
		writeLog(info, "DEBUG")
	}
}

// LogInfo will log with level Info if the logging level is set as such
func LogInfo(info datatypes.Logging) {
	level := getLoggingLevel()
	if level == "INFO" || level == "DEBUG" {
		writeLog(info, "INFO")
	}
}

// getLoggingLevel returns the logging level the app is set to
func getLoggingLevel() string {
	return logLevel
}

// writeLog prepares the log for printing
func writeLog(info datatypes.Logging, level string) {
	logMsg := "{"

	// Add level
	logMsg += fmt.Sprintf(`"level":"%v",`, level)

	// Add message
	logMsg += fmt.Sprintf(`"message":"%v",`, info.Message)

	// Possibly add perfSig
	if len(info.PerfSig.Metrics) > 0 {
		sig, _ := json.Marshal(info.PerfSig)
		logMsg += string(sig)
	}

	// Add source logging method
	file, line, method := traceMethod()
	source := fmt.Sprintf(`"source":"%v#%v:%v"`, shortenedFileName(file), shortenedMethodName(method), line)
	logMsg += source

	log.Println(logMsg + "}")
}

func shortenedFileName(name string) string {
	return strings.TrimPrefix(name, "/go/src/github.com/barrebre/goDynaPerfSignature/")
}

func shortenedMethodName(name string) string {
	method := strings.TrimPrefix(name, "github.com/barrebre/goDynaPerfSignature/")
	methodPOS := strings.LastIndex(method, ".") + 1
	return method[methodPOS:]
}

// traceMethod gives us the current method and line number. Useful for logging
// https://stackoverflow.com/questions/25927660/how-to-get-the-current-function-name
func traceMethod() (string, int, string) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	return file, line, fn.Name()
}
