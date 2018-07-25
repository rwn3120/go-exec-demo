package main

type Get struct {
    id           string
    key          string
}

func (g *Get) CorrelationId() string {
    return g.id
}

type GetResult struct {
    id    string
    value string
    error error
}

func (gr *GetResult) CorrelationId() string {
    return gr.id
}

func (gr *GetResult) Error() error {
    return gr.error
}

type Put struct {
    id     string
    key    string
    value  string
}

func (p *Put) CorrelationId() string {
    return p.id
}