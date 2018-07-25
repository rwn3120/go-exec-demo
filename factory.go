package main

import (
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-exec"
)

type Factory struct {
    configuration *Configuration
}

func (f *Factory) Handler(uuid string) exec.Handler {
    return &Handler{
        uuid:          uuid,
        configuration: f.configuration,
        logger:        logger.New(uuid+"-handler", f.configuration.Logger)}
}

func NewFactory(configuration *Configuration) exec.Factory {
    return &Factory{configuration: configuration}
}
