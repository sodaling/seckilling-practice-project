package models

type Product struct {
	ID           int64  `json:"ID" form:"ID" sql:"ID"`
	ProductName  string `json:"productName" form:"ProductName" sql:"productName"`
	ProductNum   int64  `json:"productNum" form:"ProductNum" sql:"productNum"`
	ProductImage string `json:"productImage" form:"ProductImage" sql:"productImage"`
	ProductUrl   string `json:"productUrl" form:"ProductUrl" sql:"productUrl"`
}
