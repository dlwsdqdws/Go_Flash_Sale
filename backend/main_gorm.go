package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
	"pro-iris/backend/web/controllers"
	"pro-iris/common"
	"pro-iris/repositories"
	"pro-iris/services"
)

func main() {
	// 1. Create iris instance
	app := iris.New()
	// 2. Set error mode
	app.Logger().SetLevel("debug")
	// 3. Register model
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. Set model Repository
	app.StaticWeb("/assets", "./backend/web/assets")
	// 5. Error handler
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Error Occurred！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 6. Connect database
	db, err := common.NewGormConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 7. Register controller and routing
	productRepository := repositories.NewGormProductManager(db)
	productService := services.NewGormProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.GormProductController))

	orderRepository := repositories.NewGormOrderManager(db)
	orderService := services.NewGormOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.GormOrderController))
	// 8. Start
	app.Run(
		iris.Addr("localhost:8086"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
