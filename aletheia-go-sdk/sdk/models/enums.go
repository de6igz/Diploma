package models

type LogLevel string

const (
	LogLevelDebug    LogLevel = "DEBUG"
	LogLevelInfo     LogLevel = "INFO"
	LogLevelWarning  LogLevel = "WARNING"
	LogLevelError    LogLevel = "ERROR"
	LogLevelCritical LogLevel = "CRITICAL"
)

type EventType string

const (
	EventError         EventType = "ERROR"
	EventMessage       EventType = "MESSAGE"
	EventResourceUsage EventType = "RESOURCE_USAGE"
)
