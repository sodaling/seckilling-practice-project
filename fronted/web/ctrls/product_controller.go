package ctrls

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
	"strconv"
)

var (
	//生成的Html保存目录
	htmlOutPath = "./fronted/web/htmlProductShow/"
	//静态文件模版目录
	templatePath = "./fronted/web/views/template/"
)

func GenerateHtml(c *gin.Context) {
	productString := c.Param("productID")
	productId, err := strconv.Atoi(productString)
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productSev := services.DefaultProductService()
	product, err := productSev.GetProductByID(int64(productId))
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	contentTemp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	generateStaticHtml(contentTemp, fileName, product)
}

func generateStaticHtml(template *template.Template, fileName string, product *models.Product) {
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			log.Println(err)
		}
	}
	err := os.MkdirAll(filepath.Dir(fileName), 0777)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	template.Execute(file, product)
}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
