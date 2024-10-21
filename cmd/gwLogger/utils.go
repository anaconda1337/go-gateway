package gwLogger

import (
	"fmt"
	"strings"
)

func parseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return INFO, fmt.Errorf("invalid gwLogger level: %s", level)
	}
}

func getColorByLogLevel(level LogLevel) string {
	switch level {
	case DEBUG:
		return Blue
	case INFO:
		return Green
	case WARN:
		return Yellow
	case ERROR:
		return Red
	default:
		return Reset
	}
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
