package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logWriteMu sync.Mutex

const (
	logLevelInfo  = "INFO"
	logLevelWarn  = "WARN"
	logLevelError = "ERR"
	logLevelDebug = "DEBUG"
)

func (a *Agent) writeLog(message string) {
	a.writeLogWithLevel(logLevelInfo, message)
}

func (a *Agent) writeError(message string) {
	a.writeLogWithLevel(logLevelError, message)
}

func (a *Agent) writeWarn(message string) {
	a.writeLogWithLevel(logLevelWarn, message)
}

func (a *Agent) writeDebug(message string) {
	a.writeLogWithLevel(logLevelDebug, message)
}

func (a *Agent) writeLogWithLevel(level, message string) {
	logWriteMu.Lock()
	defer logWriteMu.Unlock()

	level = strings.TrimSpace(strings.ToUpper(level))
	if level == "" {
		level = logLevelInfo
	}

	logDir := filepath.Join(os.TempDir(), "SysTrace_Agent")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("writeLog: could not create log directory: %v\n", err)
		return
	}

	logFilePath := filepath.Join(logDir, "agent.log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("writeLog: could not open log file: %v\n", err)
		return
	}
	defer file.Close()

	entry := fmt.Sprintf("%s | %s | %s\n", time.Now().Format(time.RFC3339), level, message)
	if _, err := file.WriteString(entry); err != nil {
		fmt.Printf("writeLog: could not write to log file: %v\n", err)
	}
}
