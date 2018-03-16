package phttp

type appliable interface {
    Apply(*Context, []appliable, int)
}

//appliable实现1：普通路由
type handler func(*Context) error
func (h handler) Apply(ctx *Context, apps []appliable, current int) {
    err := h(ctx)
    if err != nil {
        //TODO
    }

    //必定执行下一个appliable
    current++
    if len(apps) > current {
        apps[current].Apply(ctx, apps, current)
    }
}

//appliable实现2：通用拦截器
type middleware func(*Context, Next)
type Next func()
func (m middleware) Apply(ctx *Context, apps []appliable, current int) {
    next := func() {
        current++
        if len(apps) > current {
            apps[current].Apply(ctx, apps, current)
        }
    }
    //执行拦截逻辑
    //由拦截器处理函数决定是否执行下一个applicable(next是否被调用)
    m(ctx, next)
}


