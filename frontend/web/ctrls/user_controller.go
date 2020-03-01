package ctrls

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
	"strconv"
)

func UserRegisterShowAction(c *gin.Context) {
	c.HTML(http.StatusOK, "msg/register.html", nil)
}

func UserLoginShowAction(c *gin.Context) {
	c.HTML(http.StatusOK, "msg/login.html", nil)
}

func UserLoginAction(c *gin.Context) {
	var (
		userName = c.PostForm("userName")
		password = c.PostForm("password")
	)
	userSer := services.DefaultUserService()
	user, isOK := userSer.IsPwdSuccess(userName, password)
	if !isOK {
		log.Println(userName + ":password is wrong")
		c.Redirect(http.StatusSeeOther, "/user/login")
		return
	}
	login(user, c)
	c.Redirect(http.StatusSeeOther, "/product")
}

func login(user *models.User, c *gin.Context) {
	uidInt64 := strconv.FormatInt(user.ID, 10)
	c.SetCookie("uid", uidInt64, 3600*24, "/", "localhost", false, true)
	uidByte := []byte(uidInt64)
	uidString, err := common.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}

	c.SetCookie("sign", uidString, 3600*24, "/", "localhost", false, true)
}

func UserRegisterAction(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBind(user); err != nil {
		log.Panicln(err)
		c.Redirect(http.StatusSeeOther, "/user/error")
		return
	}
	userSer := services.DefaultUserService()
	_, err := userSer.AddUser(user)
	if err != nil {
		log.Panicln(err)
		c.Redirect(http.StatusSeeOther, "/user/error")
		return
	}
	c.Redirect(http.StatusSeeOther, "/user/login")
}
