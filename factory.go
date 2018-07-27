package main

import (
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-exec"
    "fmt"
)

type Factory struct {
    configuration *Configuration
    counter int
}

func (f *Factory) Processor() exec.Processor {
    f.counter++
    uuid := fmt.Sprintf("%s-processor-%d", f.configuration.Name, f.counter)
    return &Processor{
        configuration: f.configuration,
        logger:        logger.New(uuid      , f.configuration.Logger)}
}

func NewFactory(configuration *Configuration) exec.Factory {
    return &Factory{configuration: configuration}
}
