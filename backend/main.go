package main

import (
	"github.com/gin-gonic/gin"
	"seckilling-practice-project/backend/web/ctrls"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("backend/web/views/**/*")
	router.Static("/assets", "backend/web/assets")
	router.GET("/order", ctrls.OrderShowAction)

	router.Run(":8000")
}
