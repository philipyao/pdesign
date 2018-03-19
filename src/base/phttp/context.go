package phttp

import (
    "net/http"
)

type Context struct {
    req *Request
    rsp *Response

    sess *SessionStore
}

func makeContext(w http.ResponseWriter, r *http.Request) *Context {
    //todo r parseform
    return &Context{
        req: makeRequest(r),
        rsp: makeResponse(w, r),
    }
}

func (c *Context) Request() *Request {
    return c.req
}

func (c *Context) Response() *Response {
    return c.rsp
}

func (c *Context) Session() *SessionStore {
    return c.sess
}