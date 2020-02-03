package models

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"ID"`
	NickName     string `json:"nickname" form:"nickName" sql:"nickName"`
	UserName     string `json:"username" form:"userName" sql:"userName"`
	HashPassword string `json:"hashpassword" form:"password" sql:"passWord"`
}
