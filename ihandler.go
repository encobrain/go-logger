package logger

type IHandler interface {
	UsedFields () []string
	Handle (log *Log, fields map[string]interface{}) (final bool)
}
