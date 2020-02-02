package ctrls

import (
	"github.com/gin-gonic/gin"
	"log"
	"seckilling-practice-project/services"
)

var defaultOrderSer = services.DefaultOrderService()

func OrderShowAction(c *gin.Context) {
	orderArray, err := defaultOrderSer.GetAllOrderInfo()
	if err != nil {
		log.Fatalln("查询订单信息失败", err)
	}
	c.HTML(200, "msg/order_view.html", gin.H{"order": orderArray,})
}
