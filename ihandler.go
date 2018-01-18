package logger

type IHandler interface {
	Fields () map[string]string
	Final () bool
	Handle (log *Log, fields map[string]interface{})
}
