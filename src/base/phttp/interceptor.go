package phttp

type interceptor struct {
    middlewares []appliable
}

func (itr *interceptor) Use(mw middleware) {
    itr.middlewares = append(itr.middlewares, mw)
}

