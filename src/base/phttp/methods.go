package phttp

import (
    "fmt"
)

const (
    GET         = "GET"
    POST        = "POST"
)

type routemap map[string][]*route
type methods struct {
    rmap routemap
}
func (m *methods) initMethod() {
    m.rmap = routemap{
        GET:    []*route{},
        POST:   []*route{},
    }
    fmt.Printf("m %p, initMethod map %+v\n", m, m.rmap)
}

func (m *methods) routesBy(method string) []*route {
    return m.rmap[method]
}

func (m *methods) register(path, method string, fn handler) {
    routes, ok := m.rmap[method]
    if !ok {
        panic(fmt.Sprintf("unsupported method %v", method))
    }
    routes = append(routes, &route{
        path: path,
        fn: fn,
    })
    m.rmap[method] = routes
    fmt.Printf("m %p, register <%v>, map %+v\n", m, method, m.rmap)
}

func (m *methods) Get(path string, fn handler) {
    m.register(path, GET, fn)
}

func (m *methods) Post(path string, fn handler) {
    m.register(path, POST, fn)
}
