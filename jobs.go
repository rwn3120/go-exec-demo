package main

type Get struct {
    key          string
}

type GetResult struct {
    value string
    error error
}

func (gr *GetResult) Err() error {
    return gr.error
}

type Put struct {
    key    string
    value  string
}
