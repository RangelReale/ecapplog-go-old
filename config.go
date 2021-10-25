package ecapplog

type Priority string

const (
	Priority_TRACE = "TRACE"
	Priority_DEBUG = "DEBUG"
	Priority_INFORMATION = "INFORMATION"
	Priority_NOTICE = "NOTICE"
	Priority_WARNING = "WARNING"
	Priority_FATAL = "FATAL"
	Priority_CRITICAL = "CRITICAL"
	Priority_ERROR = "ERROR"
)
