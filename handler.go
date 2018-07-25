package main

import (
    "errors"
    "fmt"
    "time"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-exec"
    "sync"
)

var lock = sync.RWMutex{}
var store = make(map[string]string)

type Handler struct {
    uuid          string
    client        bool
    configuration *Configuration
    logger        *logger.Logger
}

func (h *Handler) Initialize() error {
    h.logger.Trace("Connecting to store")
    h.client = true
    return nil
}

func (h *Handler) Destroy() {
    h.logger.Trace("Disconnecting from store")
}

func (h *Handler) get(get *Get) exec.Result {
    var err error
    var value string
    lock.Lock()
    _, ok := store[get.key]
    lock.Unlock()
    if ok {
        <-time.After(123 * time.Millisecond) // it takes some time to get record
        lock.Lock()
        value = store[get.key]
        lock.Unlock()
    } else {
        err = errors.New(fmt.Sprintf("Key %s has not bee found", get.key))
    }

    return &GetResult{get.id, value, err}
}

func (h *Handler) put(put *Put) exec.Result {
    <-time.After(456 * time.Millisecond) // it takes some time to put record
    lock.Lock()
    store[put.key] = put.value
    lock.Unlock()
    return exec.NewResult(put.id, nil)
}

func (h *Handler) Handle(job exec.Job) exec.Result {
    if !h.client {
        panic("Handler has not been initialized")
    }

    h.logger.Trace("Processing job %s", job.CorrelationId())
    switch j := job.(type) {
    case *Put:
        return h.put(j)
    case *Get:
        return h.get(j)
    default:
        err := errors.New(fmt.Sprintf("unknown job %s", job.CorrelationId()))
        h.logger.Error("Error: %s", err)
        return exec.NewResult(job.CorrelationId(), err)
    }
}
