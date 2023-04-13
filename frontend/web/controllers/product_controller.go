package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"pro-iris/datamodels"
	"pro-iris/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
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

func (p *ProductController) GetOrder() mvc.View {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "Buy Failed!"
	// Check product number
	if product.ProductNum > 0 {
		// deduct product number
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		// create an order
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "Buy Successfully!"
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}
