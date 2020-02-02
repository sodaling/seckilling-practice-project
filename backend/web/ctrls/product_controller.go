package ctrls

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
)

var defaultProSer = services.DefaultProductService()

func ProductListAction(c *gin.Context) {

	productArray, _ := defaultProSer.GetAllProduct()
	c.HTML(200, "msg/product_view.html", gin.H{"productArray": productArray,})
}

func ProductUpdateAction(c *gin.Context) {
	product := &models.Product{}
	if err := c.ShouldBind(product);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	}
	fmt.Println(*product)
	err := defaultProSer.UpdateProduct(product)
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	}
	c.Redirect(http.StatusTemporaryRedirect,"/product/")
}