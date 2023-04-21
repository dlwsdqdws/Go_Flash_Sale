package middleware

import (
	"github.com/kataras/iris"
	"golang.org/x/time/rate"
	"time"
)

var LimiterMap = map[string]*rate.Limiter{
	"/product": rate.NewLimiter(rate.Every(time.Millisecond), 10000),
}

func TokenLimiter(ctx iris.Context) {
	var (
		uri = ctx.Request().RequestURI
	)
	limiter, check := LimiterMap[uri]
	if check {
		if allow := limiter.Allow(); !allow {
			ctx.Application().Logger().Debug("Please try again later")
			ctx.Redirect("/product")
			return
		}
	}
	ctx.Application().Logger().Debug("Successfully entered product detail's page")
	ctx.Next()
}
