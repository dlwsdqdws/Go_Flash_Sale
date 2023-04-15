package controllers

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"html/template"
	"os"
	"path/filepath"
	"pro-iris/datamodels"
	"pro-iris/rabbitmq"
	"pro-iris/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

var (
	// Directory for saving the generated HTML
	htmlOutPath = "./frontend/web/htmlProductShow/"
	// Directory for static file templates.
	templatePath = "./frontend/web/views/template/"
)

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	// 1. Check if the static file exists and delete it if already exist.
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}
	// 2.Generate a static file.
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product)
}

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	message := datamodels.NewMessage(userID, productID)
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")
}
