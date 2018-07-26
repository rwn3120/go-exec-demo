package main

import (
    "errors"
    "fmt"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-exec"
)

type Processor struct {
    uuid          string
    storageClient *StorageMockupClient
    configuration *Configuration
    logger        *logger.Logger
}

func (p *Processor) Initialize() error {
    p.logger.Trace("Connecting to store")
    if p.storageClient == nil {
        p.storageClient = &StorageMockupClient{}
    }
    p.storageClient.Connect()
    return nil
}

func (p *Processor) Destroy() {
    p.logger.Trace("Disconnecting from store")
    p.storageClient.Disconnect()
}

func (p *Processor) Process(payload exec.Payload) exec.Result {
    p.logger.Trace("Processing payload")
    switch j := payload.(type) {
    case *Put:
        err := p.storageClient.Put(j.key, j.value)
        return exec.NewResult(err)
    case *Get:
        value, err := p.storageClient.Get(j.key)
        return &GetResult{value, err}
    default:
        return exec.NewResult(errors.New(fmt.Sprintf("unknown payload")))
    }
}
