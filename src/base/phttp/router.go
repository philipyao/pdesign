package phttp

//单路由
type route struct {
    //路由path，精确匹配
    path string
    //路由处理方法
    fn  handler
    //最终路由处理链
    fnChain func(interface{})
}

//路由组
type RouteGroup struct {
    *interceptor
    *methods
}

//路由管理器
type router struct {
    *interceptor
    *methods
    groups []*RouteGroup
}

//新建路由组
func (r *router) NewGroup() *RouteGroup {
    g := &RouteGroup{}
    g.initMethod()
    r.groups = append(r.groups, g)
    return g
}

//把全局的interceptor和group里的interceptor、routemap一起merge到全局routemap里
func (r *router) mergeRoute() {
    //首先合并全局interceptor
    for _, routes := range r.rmap {
        for _, route := range routes {
            apps := merge(r.middlewares, []appliable{route.fn})
            route.fnChain = routeMakeHandlers(apps)
        }
    }

    //再合并group里的interceptor和routemap到全局routemap
    for _, group := range r.groups {
        for method, routes := range group.rmap {
            for _, route := range routes {
                apps := merge(r.middlewares, group.middlewares, []appliable{route.fn})
                route.fnChain = routeMakeHandlers(apps)
                //添加到全局routemap里
                r.rmap[method] = append(r.rmap[method], route)
            }
        }
    }
    //not need groups anymore
    r.groups = []*RouteGroup{}
}

func merge(apps ...[]appliable) []appliable {
    result := make([]appliable, 0)
    for _, app := range apps {
        result = append(result, app...)
    }
    return result
}

func routeMakeHandlers(apps []appliable) func(interface{}) {
    return func(ctx interface{}) {
        current := 0
        apps[0].Apply(ctx, apps, current)
    }
}