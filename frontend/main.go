package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
	"pro-iris/common"
	"pro-iris/frontend/middleware"
	"pro-iris/frontend/web/controllers"
	"pro-iris/rabbitmq"
	"pro-iris/repositories"
	"pro-iris/services"
	"pro-iris/tool"
	"time"
)

func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Next()
}

func main() {
	// 1. Create iris instance
	app := iris.New()
	app.Use(Cors)
	// 2. Set error mode
	app.Logger().SetLevel("debug")
	// 3. Register model
	template := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4. Set model Repository
	app.StaticWeb("/public", "./frontend/web/public")

	startTime := time.Date(2023, 4, 21, 12, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 4, 28, 12, 0, 0, 0, time.UTC)
	app.Get("/html/htmlProduct.html", middleware.OnlyDuringMiddleware(startTime, endTime), func(ctx iris.Context) {
		ctx.ServeFile("./frontend/web/htmlProductShow/htmlProduct.html", false)
	})
	app.StaticWeb("/html", "./frontend/web/htmlProductShow")

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
	// 7. Register controller and routing
	user := repositories.NewUserManagerRepository("user", db)
	userService := services.NewUserService(user)
	proUser := app.Party("/user")
	userParty := mvc.New(proUser)
	userParty.Register(userService, ctx)
	userParty.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("rabbitmqProduct")

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	productParty := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	proProduct.Use(middleware.TokenLimiter)
	productParty.Register(productService, orderService, ctx, rabbitmq)
	productParty.Handle(new(controllers.ProductController))

	app.Get("/captcha/", tool.GetCaptchaID)
	app.Get("/captcha/*", tool.GetCaptchaImg)
	app.Get("/captcha/verify", tool.VerifyCaptcha)

	// 8. Start
	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
