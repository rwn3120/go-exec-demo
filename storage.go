package main

import (
    "strings"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-conf"
    "github.com/rwn3120/go-exec"
    "github.com/satori/go.uuid"
)

type Configuration struct {
    Name               string
    Logger             *logger.Configuration
    Executor           *exec.Configuration
}

func (c Configuration) Validate() *[]string {
    var errorList []string

    if len(strings.TrimSpace(c.Name)) == 0 {
        errorList = append(errorList, "Configuration: Missing name")
    }

    otherErrors := conf.Validate(c.Logger, c.Executor)
    if otherErrors != nil {
        errorList = append(errorList, *otherErrors...)
    }

    if errorsCount := len(errorList); errorsCount > 0 {
        return &errorList
    }
    return nil
}

type Storage struct {
    configuration *Configuration
    executor      *exec.Executor
    logger        *logger.Logger
}

func New(configuration *Configuration) *Storage {
    // check configuration
    conf.Check(configuration)

    // create backend
    return &Storage{
        configuration: configuration,
        executor: exec.New(
            configuration.Name,
            configuration.Executor,
            NewFactory(configuration)),
        logger: logger.New(
            "storage",
            configuration.Logger)}
}

func (b *Storage) Destroy() {
    b.executor.Destroy()
}

func (b *Storage) Put(key string, value string) error {
    channel := make(chan exec.Result)
    callback := func(result exec.Result) {
        channel <- result
    }
    b.executor.Fire(&Put{
        id:    uuid.NewV4().String(),
        key:   key,
        value: value},
        callback)
    result := <-channel
    return result.Error()
}

func (b *Storage) Get(key string) (string, error) {
    result := b.executor.Execute(&Get{
        id:  uuid.NewV4().String(),
        key: key}).(*GetResult)
    return result.value, result.error
}
