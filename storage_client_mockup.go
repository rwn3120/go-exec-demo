package main

import (
    "sync"
    "time"
    "fmt"
    "errors"
)

var lock = sync.RWMutex{}
var store = make(map[string]string)

type StorageClientMockup struct {
    connected bool
}

func (p *StorageClientMockup) Connect() {
    p.connected = true
}

func (p *StorageClientMockup) Disconnect() {
    p.connected = false
}

func (p *StorageClientMockup) Connected() bool {
    return p.connected
}

func (p *StorageClientMockup) Get(key string) (string, error) {
    if !p.Connected(){
        return "", errors.New("connection error")
    }

    lock.Lock()
    _, ok := store[key]
    lock.Unlock()
    if ok {
        <-time.After(123 * time.Millisecond) // it takes some time to get record
        lock.Lock()
        value := store[key]
        lock.Unlock()
        return value, nil
    } else {
        return "", errors.New(fmt.Sprintf("Key %s has not bee found", key))
    }
}

func (p StorageClientMockup) Put(key string, value string) error {
    if !p.Connected(){
        return errors.New("connection error")
    }

    <-time.After(456 * time.Millisecond) // it takes some time to put record
    lock.Lock()
    store[key] = value
    lock.Unlock()
    return nil
}
