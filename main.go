package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"pro-iris/web/controllers"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("./web/views", ".html"))
	//register controller
	mvc.New(app.Party("/hello")).Handle(new(controllers.MovieController))
	app.Run(
		iris.Addr("localhost:8080"),
	)
}
