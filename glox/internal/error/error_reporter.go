package error

type ErrorReporter interface {
	Error(line int, message string)
	InternalError(message string)
}
