package main

import "github.com/kataras/iris/v12"

func main() {
	// 1. Create iris instance
	app := iris.New()
	// 2. Set error mode
	app.Logger().SetLevel("debug")
	// 3. Register model
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. Set model Repository
	app.HandleDir("/assets", iris.Dir("./backend/web/assets"))
	// 5. Error handler
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Error Occurred!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 6. Register controller and routing
	// 7. Start
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
