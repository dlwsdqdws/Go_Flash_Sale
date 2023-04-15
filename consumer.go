package main

import (
	"fmt"
	"pro-iris/common"
	"pro-iris/rabbitmq"
	"pro-iris/repositories"
	"pro-iris/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	rabbitmqConsumerSimple := rabbitmq.NewRabbitMQSimple("rabbitmqProduct")
	rabbitmqConsumerSimple.ConsumeSimple(orderService, productService)
}
