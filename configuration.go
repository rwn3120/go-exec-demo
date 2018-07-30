package main

import (
    "strings"
    "github.com/rwn3120/go-conf"
    "github.com/rwn3120/go-logger"
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