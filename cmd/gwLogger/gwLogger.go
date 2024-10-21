package gwLogger

import (
	"fmt"
	"go-gateway/cmd/conf"
	"os"
	"runtime"
	"time"
)

func NewLogger(config *conf.Config) (*Logger, error) {
	level, err := parseLogLevel(config.GatewayConfig.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		logLevel:  level,
		logToFile: config.GatewayConfig.LogToFile,
	}

	if config.GatewayConfig.LogToFile {
		file, err := os.OpenFile(config.GatewayConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open gwLogger file: %v", err)
		}
		logger.logFile = file
	}

	return logger, nil
}

func (l *Logger) logMessage(level LogLevel, message string) {
	if level >= l.logLevel {
		timestamp := time.Now().Format(time.RFC3339)

		pc, file, line, ok := runtime.Caller(2)
		funcName := ""
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}

		color := getColorByLogLevel(level)

		logLine := fmt.Sprintf("%s[%s] %s [%s:%d %s] %s%s\n",
			color, levelToString(level), timestamp, file, line, funcName, message, Reset)

		if l.logToFile && l.logFile != nil {
			_, err := l.logFile.WriteString(fmt.Sprintf("[%s] %s [%s:%d %s] %s\n",
				levelToString(level), timestamp, file, line, funcName, message))
			if err != nil {
				fmt.Printf("failed to write to log file: %v", err)
			}
		} else {
			fmt.Print(logLine)
		}
	}
}

func (l *Logger) Info(message string) {
	l.logMessage(INFO, message)
}

func (l *Logger) Debug(message string) {
	l.logMessage(DEBUG, message)
}

func (l *Logger) Warn(message string) {
	l.logMessage(WARN, message)
}

func (l *Logger) Error(message string) {
	l.logMessage(ERROR, message)
}

func (l *Logger) Close() {
	if l.logToFile && l.logFile != nil {
		err := l.logFile.Close()
		if err != nil {
			return
		}
	}
}
