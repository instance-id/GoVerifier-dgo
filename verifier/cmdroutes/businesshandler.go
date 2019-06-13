package cmdroutes

import (
	"github.com/Necroforger/dgrouter/exrouter"
)

type Route interface {
	GetCommand() string
	GetDescription() string
}

type DefaultRouteHandler interface {
	Handle(ctx *exrouter.Context)
}

type CustomRouteHandler interface {
	Register(router *exrouter.Route) *exrouter.Route
}

type SubRoute interface {
	Route
	GetSubRoutes() []Route
}

func registerHelpHandler(router *exrouter.Route) {
	route := NewHelpRoute(router)
	registerRoute(route, router)
}

func registerRoute(route Route, router *exrouter.Route) *exrouter.Route {
	prefix := route.GetCommand()
	var newRouter *exrouter.Route
	if handlerHandler, ok := route.(DefaultRouteHandler); ok {
		newRouter = router.On(prefix, handlerHandler.Handle)
	}
	if routeHandler, ok := route.(CustomRouteHandler); ok {
		newRouter = routeHandler.Register(router)
	}
	newRouter.Desc(route.GetDescription())
	return newRouter
}

func RegisterRoutes(router *exrouter.Route, routes ...Route) {
	for _, route := range routes {
		registerHelpHandler(router)

		rootRouter := registerRoute(route, router)
		if subRoute, ok := route.(SubRoute); ok {
			if subRouteHandlers := subRoute.GetSubRoutes(); subRouteHandlers != nil {
				registerHelpHandler(rootRouter)

				for _, handler := range subRouteHandlers {
					registerRoute(handler, rootRouter)
				}
			}
		}
		//log.Printf("Handlers registered: %s", route.GetCommand())
	}

}
