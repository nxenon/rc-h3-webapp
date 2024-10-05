package models

type ProductObject struct {
	ProductId       int
	ProductName     string
	ProductPrice    int
	ProductImageUrl string
	ProductInCartId int
}

type AddProductRequest struct {
	ProductId int `json:"ProductId"`
}

type RemoveProductRequest struct {
	ProductInCartId int `json:"ProductInCartId"`
}
