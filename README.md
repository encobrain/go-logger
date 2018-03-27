Usage

```go
package main

import "github.com/encobrain/go-logger"

type handler struct {}

func (h *handler) UsedFields () []string {
	return []string{}
}

func (h *handler) Handle (log *logger.Log, fields map[string]interface{}) (final bool) {
    return false	
}

func main() {
	log := &logger.Log{}
	
	log = log.AddHandler(&handler{})
	
	log = log.Fields("abc", 123)
   
    log.Tracef("some log text")
    log.Debugf("some log text")
    log.Infof("some log text")
    log.Warnf("some log text")
    log.Errorf("some log text")	
    
    log.Handle()
}



```