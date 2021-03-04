package webserver

import (
	"github.com/bytepowered/flux"
	"github.com/labstack/echo/v4"
)

// RouteHandler 实现flux.WebRouteHandler与Echo框架的echo.HandlerFunc函数适配
type AdaptWebRouteHandler flux.WebHandler

func (f AdaptWebRouteHandler) AdaptFunc(ctx echo.Context) error {
	return f(newAdaptWebContext(ctx))
}

// AdaptWebInterceptor 实现flux.WebInterceptor与Echo框架的echo.MiddlewareFunc函数适配
type AdaptWebInterceptor flux.WebInterceptor

func (intfun AdaptWebInterceptor) AdaptFunc(next echo.HandlerFunc) echo.HandlerFunc {
	handler := intfun(func(webc flux.WebContext) error {
		return next(webc.(*AdaptWebContext).echoc)
	})
	return func(echoc echo.Context) error {
		return handler(newAdaptWebContext(echoc))
	}
}
