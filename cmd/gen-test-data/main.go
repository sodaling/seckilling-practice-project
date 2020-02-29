// 生成批量用于测试数据用的用户和商品
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
	"strconv"
)

func main() {
	userSer := services.DefaultUserService()
	productSer := services.DefaultProductService()

	var users []*models.User
	var products []*models.Product

	// 生成1000个用户，用户名格式为soda1到soda1000,密码都为1234
	for i := 1; i <= 1000; i++ {
		uName := "soda" + strconv.Itoa(i)
		user := &models.User{UserName: uName, NickName: uName, HashPassword: "1234"}
		uid, err := userSer.AddUser(user)
		if err != nil {
			fmt.Println(err)
		}
		user.ID = uid
		fmt.Println(user)
		fmt.Println(uName, " has been created")
		users = append(users, user)
	}

	// 生成10个商品，商品名等等都为green1到green10
	for i := 1; i <= 5; i++ {
		pName := "soda" + strconv.Itoa(i)
		product := &models.Product{ProductName: pName, ProductNum: 100000, ProductImage: "image", ProductUrl: "url"}
		_, err := productSer.InsertProduct(product)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(pName, " has been created")
		products = append(products, product)
	}
	// 写入文件已备后面测试用

	f, err := os.Create("test.txt")
	defer f.Close()
	if err != nil {
		log.Println("test.txt cant be created:", err)
	}
	w := bufio.NewWriter(f)
	for _, user := range users {
		uid := strconv.FormatInt(user.ID, 10)
		uidByte := []byte(uid)
		sign, err := common.EnPwdCode(uidByte)
		if err != nil {
			fmt.Println(user, "sign err", err)
			continue
		}
		_, err = w.WriteString(uid + "," + sign + "\n")
	}
	w.Flush()

	log.Println("all job has done")
}
