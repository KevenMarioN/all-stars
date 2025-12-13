package adapter

import "context"

type ILogger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	WithError(error) ILogger
	WithContext(context.Context) ILogger
	AddField(string, any) ILogger
	WithFields(map[string]any) ILogger
}
