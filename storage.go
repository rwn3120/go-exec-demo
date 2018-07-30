package main

import (
    "github.com/rwn3120/go-exec"
    "errors"
)

type Session struct {
    uuid string
}

type Storage struct {
    session  *Session
    executor *exec.Executor
}

func (s *Storage) Put(key, value string) error {
    channel := make(chan exec.Result)
    callback := func(result exec.Result) {
        channel <- result
    }
    s.executor.FireJob(&Put{
        session: s.session,
        key:     key,
        value:   value},
        callback)
    result := <-channel
    return result.Err()
}

func (s *Storage) Get(key string) (string, error) {
    get := &Get{session: s.session, key: key}
    result, err := s.executor.ExecuteJob(get)
    if err != nil {
        return "", err
    }
    getResult, ok := result.(*GetResult)
    if !ok {
        return "", errors.New("unexpected result")
    }
    return getResult.value, getResult.error
}
