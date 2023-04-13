package middleware

import "github.com/kataras/iris"

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("Please Login Your Account")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("You have logged in")
	ctx.Next()
}
