package logger

import (
	"time"
	"fmt"
)

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

func (l *Log) Trace (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "trace", "_message", fmt.Sprintf(format, args)).Handle()
}

func (l *Log) Debug (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "debug", "_message", fmt.Sprintf(format, args)).Handle()
}

func (l *Log) Info (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "info", "_message", fmt.Sprintf(format, args)).Handle()
}

func (l *Log) Warn (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "warn", "_message", fmt.Sprintf(format, args)).Handle()
}

func (l *Log) Error (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "error", "_message", fmt.Sprintf(format, args)).Handle()
}

func (l *Log) Panic (format string, args ...interface{}) {
	l.Fields("_datetime", time.Now(), "_level", "panic", "_message", fmt.Sprintf(format, args)).Handle()
}


func (l *Log) Handle () {
	fbits := l.fbits
	fields := l.fields

	handlers := l.fbitsCache[fbits]

	if handlers == nil {
		next:
		for _,h := range l.handlers {
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

		l.fbitsCache[l.fbits] = handlers
	}

	if len(handlers) == 0 {
		fmt.Printf("%#v\n", fields)
	} else {
		for _,h := range handlers {
			h.Handle(l, fields)
		}
	}


}


