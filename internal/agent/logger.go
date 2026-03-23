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

	entry := fmt.Sprintf("%s | %s | %s\n", time.Now().Format(time.RFC3339), level, message)

	var lastErr error
	for _, logDir := range resolveLogDirCandidates() {
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			lastErr = fmt.Errorf("create log directory %s: %w", logDir, err)
			continue
		}

		logFilePath := filepath.Join(logDir, "agent.log")
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			lastErr = fmt.Errorf("open log file %s: %w", logFilePath, err)
			continue
		}

		if _, err := file.WriteString(entry); err != nil {
			_ = file.Close()
			lastErr = fmt.Errorf("write log file %s: %w", logFilePath, err)
			continue
		}

		_ = file.Close()
		return
	}

	if lastErr != nil {
		fmt.Printf("writeLog: could not write log entry: %v\n", lastErr)
	}
}

func resolveLogDirCandidates() []string {
	paths := make([]string, 0, 3)

	if wd, err := os.Getwd(); err == nil && wd != "" {
		paths = append(paths, filepath.Join(wd, "logs"))
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		if exeDir != "" {
			paths = append(paths, filepath.Join(exeDir, "logs"))
		}
	}

	paths = append(paths, filepath.Join(os.TempDir(), "SysTrace_Agent"))
	return paths
}
