package main

import (
    "errors"
    "fmt"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-exec"
)

type Processor struct {
    uuid          string
    storageClient *StorageClientMockup
    configuration *Configuration
    logger        *logger.Logger
}

func (p *Processor) Initialize() error {
    p.logger.Trace("Connecting to store")
    if p.storageClient == nil {
        p.storageClient = &StorageClientMockup{}
    }
    p.storageClient.Connect()
    return nil
}

func (p *Processor) Destroy() {
    p.logger.Trace("Disconnecting from store")
    p.storageClient.Disconnect()
}

func (p *Processor) Process(payload exec.Payload) exec.Result {
    switch payload := payload.(type) {
    case *Put:
        p.logger.Debug("[%s] Put: %10s <- %s", payload.session.uuid, payload.key, payload.value)
        err := p.storageClient.Put(payload.key, payload.value)
        return exec.NewResult(err)
    case *Get:
        p.logger.Debug("[%s] Get: %10s", payload.session.uuid, payload.key)
        value, err := p.storageClient.Get(payload.key)
        return &GetResult{value, err}
    default:
        return exec.Nok(errors.New(fmt.Sprintf("unsupoorted payload")))
    }
}

type ProcessorFactory struct {
    configuration *Configuration
    counter       int
}

func (f *ProcessorFactory) Processor() exec.Processor {
    f.counter++
    uuid := fmt.Sprintf("%s-processor-%d", f.configuration.Name, f.counter)
    return &Processor{
        configuration: f.configuration,
        logger:        logger.New(uuid, f.configuration.Logger)}
}
