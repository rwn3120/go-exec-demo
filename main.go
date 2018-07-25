package main

import (
    "flag"
    "strings"
    "github.com/rwn3120/go-conf"
    "fmt"
    "time"
)

func main() {
    var configFile = flag.String("config", "", "path to config file")
    flag.Parse()

    if strings.TrimSpace(*configFile) == "" {
        panic("Wrong arguments")
    }
    configuration := &Configuration{}
    conf.LoadYamlAndCheck(*configFile, configuration)
    conf.PrintYaml(configuration)

    storage := New(configuration)
    type res struct {
        status bool
        key    string
    }

    count := 10000
    channel := make(chan res, count)
    startTime := time.Now()
    for i := 0; i < count; i++ {
        go func(i int) {
            key := fmt.Sprintf("key-%08d", i)
            value := fmt.Sprintf("%d", i)
            storage.Put(key, value)
            storedValue, _ := storage.Get(key)
            if value == storedValue {
                channel <- res{true, key}
            } else {
                channel <- res{false, key}
            }
        }(i)
    }

    for i := 0; i < count; i++ {
        result := <-channel
        if result.status {
            fmt.Println("OK  ", result.key)
        } else {
            fmt.Println("NOK ", result.key)
        }
    }
    fmt.Printf("Done [%v]\n", time.Since(startTime))
    storage.Destroy()
}
