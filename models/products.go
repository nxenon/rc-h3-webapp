package models

type ProductObject struct {
	ProductId         int
	ProductName       string
	ProductPrice      int
	ProductImageUrl   string
	ProductInCartId   int
	ProductInCartUUID string
}

type AddProductRequest struct {
	ProductId int `json:"ProductId"`
}

type RemoveProductRequest struct {
	ProductInCartUUID string `json:"ProductInCartUUID"`
}
