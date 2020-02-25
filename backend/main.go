package main

import (
	"github.com/gin-gonic/gin"
	"seckilling-practice-project/backend/web/ctrls"
	"seckilling-practice-project/common"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("backend/web/views/**/*")
	router.Static("/assets", "backend/web/assets")
	orderRou := router.Group("/order")
	{
		orderRou.GET("/", ctrls.OrderShowAction)
	}
	productRou := router.Group("/product")
	{
		productRou.GET("/", ctrls.ProductListAction)
		productRou.POST("/update", ctrls.ProductUpdateAction)
		productRou.GET("/manager", ctrls.ProductManagerAction)
		productRou.GET("/delete", ctrls.ProductDeleteAction)
		productRou.GET("/add", ctrls.ProductAddShowAction)
		productRou.POST("/add", ctrls.ProductAddAction)
	}

	crt, key := common.GetGrpcCrtKey()
	router.RunTLS(":8000", crt, key)
}
