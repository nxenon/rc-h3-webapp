package models

type Cart struct {
	CartId           int
	UserId           int
	CartOverallPrice int
	Products         []ProductObject
}

type MyCartResponse struct {
	UserId           int             `json:"user_id"`
	CartOverallPrice int             `json:"cart_overall_price"`
	Products         []ProductObject `json:"products"`
}
