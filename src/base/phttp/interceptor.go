package phttp

import (
    "fmt"
)

type interceptor struct {
    middlewares []appliable
}

func (itr *interceptor) Use(mw middleware) {
    itr.middlewares = append(itr.middlewares, mw)
    fmt.Printf("interceptor %p, middles count %v\n", itr, len(itr.middlewares))
}

