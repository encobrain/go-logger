package logger

import (
	"time"
	"fmt"
	"runtime"
	"strings"
)

var GOROOT string

func init() {
	_,file,_,ok := runtime.Caller(0)

	if !ok { panic("Cant get GOROOT") }

	GOROOT = strings.Replace(file, "github.com/encobrain/go-logger/log.go", "",-1)
}

type h struct {
	handler IHandler
	mask    uint64
}

type Log struct {
	fields 		map[string]interface{}
	fbits  		uint64

	handlers 	[]h
	fBitMap  	map[string]uint64

	fbitsCache 	map[uint64][]IHandler
}

func (l *Log) AddHandler (handler IHandler) *Log {
	if l.fBitMap == nil { l.fBitMap = map[string]uint64{} }

	var mask uint64

	for f := range handler.Fields() {
		bit := l.fBitMap[f]

		if bit == 0 {
			le := len(l.fBitMap)
			if le>=64 { panic("Too many used fields: >64") }

			bit = 1 << uint(le)

			l.fBitMap[f]=bit
		}

		mask |= bit
	}
	
	log := &Log{
		fields:   map[string]interface{}{},
		fbits:    l.fbits,
		handlers: append(l.handlers, h{handler: handler, mask:mask}),
		fBitMap:  l.fBitMap,

		fbitsCache: map[uint64][]IHandler{},
	}

	for f,v := range l.fields {
		log.fields[f]=v
	}

  	return log
}

func (l *Log) Fields (fields ...interface{}) *Log {
	le := len(fields)

	if le % 2 != 0 {
		panic("fields count should be even")
	}

 	log := &Log{
 		fields:     map[string]interface{}{},
 		fbits :		l.fbits,
 		handlers:	l.handlers,
 		fBitMap: 	l.fBitMap,
 		fbitsCache: l.fbitsCache,
	}

	for f,v := range l.fields {
		log.fields[f] = v
	}

	i := 0

	for i<le {
		f,ok := fields[i].(string); i++
		if !ok { panic("field name must be string") }
		
		v := fields[i]; i++
		log.fields[f]=v
		log.fbits |= log.fBitMap[f]
	}

 	return log
}

func (l *Log) Tracef (format string, args ...interface{}) {
	log := l.Fields("_level", "trace", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}

func (l *Log) Debugf (format string, args ...interface{}) {
	log := l.Fields("_level", "debug", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}

func (l *Log) Infof (format string, args ...interface{}) {
	log := l.Fields("_level", "info", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}

func (l *Log) Warnf (format string, args ...interface{}) {
	log := l.Fields("_level", "warn", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}

func (l *Log) Errorf (format string, args ...interface{}) {
	log := l.Fields("_level", "error", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}

func (l *Log) Panicf (format string, args ...interface{}) {
	log := l.Fields("_level", "panic", "_message", fmt.Sprintf(format, args...))
	_,file,line,ok := runtime.Caller(1)
	if ok { log = log.Fields("_file", file, "_line", line) }
	log.Handle()
}


func (l *Log) Handle () {
	log := l.Fields("_datetime", time.Now())

	if file,ok := log.fields["_file"]; !ok {
		_,file,line,ok := runtime.Caller(1)
		if ok {
			log = log.Fields("_file", strings.Replace(file, GOROOT, "", -1), "_line", line)
		}
	} else {
		log.fields["_file"] = strings.Replace(file.(string), GOROOT, "", -1)
	}

	fbits := log.fbits
	fields := log.fields

	handlers := log.fbitsCache[fbits]

	if handlers == nil {
		next:
		for _,h := range log.handlers {
			if h.mask & fbits == h.mask {
				handler := h.handler

				for f,v := range handler.Fields() {
					if v != "*" {
						fv,ok := fields[f].(string)

						if !ok || fv != v {
							continue next
						}
					}
				}

				handlers = append(handlers, handler)

				if handler.Final() { break }
			}
		}

		log.fbitsCache[log.fbits] = handlers
	}

	if len(handlers) == 0 {
		fmt.Printf("%#v\n", fields)
	} else {
		for _,h := range handlers {
			h.Handle(log, fields)
		}
	}


}


