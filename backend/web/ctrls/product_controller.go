package ctrls

import (
	"github.com/gin-gonic/gin"
	"seckilling-practice-project/services"
)

var defaultProSer = services.DefaultProductService()

func ProductListAction(c *gin.Context) {

	productArray, _ := defaultProSer.GetAllProduct()
	c.HTML(200, "msg/product_view.html", gin.H{"productArray": productArray,})
}
