package main

import (
    "github.com/rwn3120/go-exec"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-conf"
)

type RuntimeFactory interface {
    Storage(uuid string) (*Storage, error)
    Destroy()
}

type Runtime struct {
    configuration *Configuration
    executor      *exec.Executor
    logger        *logger.Logger
}

func (b *Runtime) Destroy() {
    b.logger.Trace("Destroying factory")
    b.executor.Destroy()
}

func (b *Runtime) Storage(uuid string) (*Storage, error) {
    b.logger.Trace("Creating executor %s", uuid)
    storage :=  &Storage{
        session:  &Session{uuid},
        executor: b.executor,
    }
    return storage, nil
}

func Production(configuration *Configuration) *Runtime {
    // check configuration
    conf.Check(configuration)
    // create executor factory
    return &Runtime{
        configuration: configuration,
        executor: exec.New(
            configuration.Name,
            configuration.Executor,
            &ProcessorFactory{configuration: configuration}),
        logger: logger.New(
            "production",
            configuration.Logger)}
}