package phttp

import (
    "net/http"
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
)

type Request struct {
    req *http.Request
}

func makeRequest(r *http.Request) *Request {
    return &Request{
        req: r,
    }
}

func (r *Request) Method() string {
    return r.req.Method
}

func (r *Request) Path() string {
    return r.req.URL.Path
}

//incoming Body as Form, get value from it
func (r *Request) FormValue(key string) string {
    return r.req.FormValue(key)
}

// Parse incoming Body as JSON
func (r *Request) JsonBody(pkg interface{}) error {
    dec := json.NewDecoder(r.req.Body)
    defer r.req.Body.Close()

    return dec.Decode(pkg)
}

// Parse incoming Body as XML
func (r *Request) XmlBody(pkg interface{}) error {
    dec := xml.NewDecoder(r.req.Body)
    defer r.req.Body.Close()

    return dec.Decode(pkg)
}

// Get raw body data
func (r *Request) RawBody() ([]byte, error) {
    content, err := ioutil.ReadAll(r.req.Body)
    if err != nil {
        return nil, err
    }
    defer r.req.Body.Close()

    return content, nil
}
