package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

var levelColors = map[LogLevel]string{
	DEBUG: colorBlue,
	INFO:  colorGreen,
	WARN:  colorYellow,
	ERROR: colorRed,
}

type CustomLogger struct {
	output io.Writer
	mu     sync.Mutex
}

func NewLogger(out io.Writer) *CustomLogger {
	return &CustomLogger{
		output: out,
	}
}

func (l *CustomLogger) logf(level LogLevel, format string, v ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	userMessage := fmt.Sprintf(format, v...)
	color, ok := levelColors[level]
	if !ok {
		color = colorReset
	}
	fullLine := fmt.Sprintf("[%s] %s%s%s: %s\n", 
		timestamp, color, levelNames[level], colorReset, userMessage)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output.Write([]byte(fullLine))
}

func (l *CustomLogger) Debugf(f string, v ...any) { l.logf(DEBUG, f, v...) }
func (l *CustomLogger) Infof(f string, v ...any)  { l.logf(INFO, f, v...) }
func (l *CustomLogger) Warnf(f string, v ...any)  { l.logf(WARN, f, v...) }
func (l *CustomLogger) Errorf(f string, v ...any) { l.logf(ERROR, f, v...) }
