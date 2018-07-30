package main

type Session struct {
    uuid string
}

type Get struct {
    session *Session
    key     string
}

type GetResult struct {
    value string
    error error
}

func (gr *GetResult) Err() error {
    return gr.error
}

type Put struct {
    session *Session
    key     string
    value   string
}