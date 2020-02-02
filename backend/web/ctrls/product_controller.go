package ctrls

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
	"strconv"
)

var defaultProSer services.IproductSerive

func init() {
	defaultProSer = services.DefaultProductService()
}

func ProductListAction(c *gin.Context) {
	productArray, _ := defaultProSer.GetAllProduct()
	c.HTML(200, "msg/product_view.html", gin.H{"productArray": productArray,})
}

func ProductUpdateAction(c *gin.Context) {
	product := &models.Product{}
	if err := c.ShouldBind(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := defaultProSer.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/product/")
}

func ProductManagerAction(c *gin.Context) {
	idString := c.Query("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := defaultProSer.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "msg/manager.html", gin.H{"product": product})
}

func ProductDeleteAction(c *gin.Context) {
	idString := c.Query("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isOk := defaultProSer.DeleteProductByID(id)
	if isOk {
		log.Println("删除商品成功，ID为：" + idString)
	} else {
		log.Panicln("删除商品失败，ID为：" + idString)
	}
	c.Redirect(http.StatusSeeOther, "/product/")
}

func ProductAddAction(c *gin.Context) {
	product := &models.Product{}
	if err := c.ShouldBind(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := defaultProSer.InsertProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Redirect(http.StatusSeeOther, "/product/")
}

func ProductAddShowAction(c *gin.Context) {
	c.HTML(http.StatusOK, "msg/add.html", nil)
}
