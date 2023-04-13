package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"github.com/opentracing/opentracing-go/log"
	"pro-iris/common"
	"pro-iris/frontend/web/controllers"
	"pro-iris/repositories"
	"pro-iris/services"
	"time"
)

func main() {
	app := iris.New()
	// 2. Set error mode
	app.Logger().SetLevel("debug")
	// 3. Register model
	template := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. Set model Repository
	app.StaticWeb("/public", "./frontend/web/public")
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")
	// 5. Error handler
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Error Occurred！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 6. Connect database
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sess := sessions.New(sessions.Config{
		Cookie:  "helloworld",
		Expires: 60 * time.Minute,
	})
	// 7. Register controller and routing
	user := repositories.NewUserManagerRepository("user", db)
	userService := services.NewUserService(user)
	userParty := mvc.New(app.Party("/user"))
	userParty.Register(userService, ctx, sess.Start)
	userParty.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	proProduct := app.Party("/product")
	productParty := mvc.New(proProduct)
	productParty.Register(productService, sess.Start)
	productParty.Handle(new(controllers.ProductController))

	// 8. Start
	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
