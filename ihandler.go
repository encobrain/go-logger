package logger

type IHandler interface {
	Fields () map[string]string
	Final () bool
	Handle (*Log, *Fields)
}
