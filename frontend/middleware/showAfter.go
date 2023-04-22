package middleware

import (
	"github.com/kataras/iris"
	"time"
)

func OnlyDuringMiddleware(startTime time.Time, endTime time.Time) iris.Handler {
	return func(ctx iris.Context) {
		if time.Now().Before(startTime) || time.Now().After(endTime) {
			//ctx.Redirect("/user/login")
			ctx.StatusCode(iris.StatusForbidden)
			ctx.WriteString("Access denied - this resource can only be accessed between " + startTime.String() + " and " + endTime.String())
			return
		}
		ctx.Next()
	}
}
