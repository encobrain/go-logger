Usage

```go
package main

import "github.com/encobrain/go-logger"

type handler struct {
	
}

func (h *handler) Fields () map[string]string {
	return map[string]string{}
}

func (h *handler) Final () bool {
	return false
}

func (h *handler) Handle (log *logger.Log, fields *logger.Fields) {
    	
}

func main() {
	log := &logger.Log{}
	
	log = log.AddHandler(&handler{})
	
	log.Fields(logger.Fields{})
   
    log.Trace("some log text")
    log.Debug("some log text")
    log.Info("some log text")
    log.Warn("some log text")
    log.Error("some log text")	
    
    log.Handle()
}



```