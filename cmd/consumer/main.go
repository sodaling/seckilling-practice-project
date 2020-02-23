package main

import (
	"fmt"
	"os"
	"seckilling-practice-project/common"
	"seckilling-practice-project/rabbitmq"
	"seckilling-practice-project/respsoiories"
	"seckilling-practice-project/services"
)

func main() {
	db, err := common.DefaultDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	productResp := respsoiories.NewProductManager("product", db)
	productSer := services.NewProductService(productResp)

	orderResp := respsoiories.NewOrderMangerRepository("order", db)
	orderSer := services.NewOrderService(orderResp)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("miaosha")
	rabbitmqConsumeSimple.ConsumeSimple(orderSer, productSer)
}
