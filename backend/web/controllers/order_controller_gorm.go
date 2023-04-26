package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"pro-iris/common"
	"pro-iris/datamodels"
	"pro-iris/services"
	"strconv"
)

type GormOrderController struct {
	Ctx          iris.Context
	OrderService services.IGormOrderService
}

func (o *GormOrderController) GetAll() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("Fail to access order information")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

func (o *GormOrderController) PostUpdate() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	err := o.OrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *GormOrderController) GetManager() mvc.View {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	order, err := o.OrderService.GetOrderByID(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "order/manager.html",
		Data: iris.Map{
			"order": order,
		},
	}
}

func (o *GormOrderController) GetDelete() {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	check := o.OrderService.DeleteOrderByID(id)
	if check {
		o.Ctx.Application().Logger().Debug("Delete successfully，ID：" + idString)
	} else {
		o.Ctx.Application().Logger().Debug("Error occurred when deleting，ID：" + idString)
	}
	o.Ctx.Redirect("/order/all")
}
