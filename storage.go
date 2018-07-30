package main

import (
    "strings"
    "github.com/rwn3120/go-logger"
    "github.com/rwn3120/go-conf"
    "github.com/rwn3120/go-exec"
)

type Configuration struct {
    Name     string
    Logger   *logger.Configuration
    Executor *exec.Configuration
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

func (b *Storage) put(session *Session, key, value string) error {
    channel := make(chan exec.Result)
    callback := func(result exec.Result) {
        channel <- result
    }
    b.executor.FireJob(&Put{
        session: session,
        key:   key,
        value: value},
        callback)
    result := <-channel
    return result.Err()
}

func (b *Storage) get(session *Session, key string) (string, error) {
    get := &Get{session: session, key: key}
    result, err := b.executor.ExecuteJob(get)
    if err != nil {
        return "", err
    }
    getResult := result.(*GetResult)
    return getResult.value, getResult.error
}

func (b *Storage) Session(uuid string) *StorageSession {
    return &StorageSession{
        session: &Session{uuid},
        storage: b,
    }
}

type StorageSession struct {
    session *Session
    storage *Storage
}

func (s *StorageSession) Put(key, value string) error {
    return s.storage.put(s.session, key, value)
}

func (s *StorageSession) Get(key string) (string, error) {
    return s.storage.get(s.session, key)
}