package models

type CartProduct struct {
	ProductId    int
	ProductPrice int
}

type UserCart struct {
	UserId           int
	CartProducts     []CartProduct
	CartOverallPrice int
}
