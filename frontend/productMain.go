package main

import (
	"github.com/kataras/iris"
)

func main() {
	// 1. Create iris instance
	app := iris.New()
	// 2. Set model Repository
	app.StaticWeb("/public", "./frontend/web/public")
	app.StaticWeb("/html", "./frontend/web/htmlProductShow")
	// 3. Start
	app.Run(
		iris.Addr("0.0.0.0:80"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
