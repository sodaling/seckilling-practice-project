package models

type Product struct {
	ID           int64  `json:"ID" form:"ID"`
	ProductName  string `json:"productName" form:"ProductName"`
	ProductNum   int64  `json:"productNum" form:"ProductNum"`
	ProductImage string `json:"productImage" form:"ProductImage"`
	ProductUrl   string `json:"productUrl" form:"ProductUrl"`
}
