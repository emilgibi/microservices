package models

type Stock struct {
	ID              int    `json:"id"`
	ProductName     string `json:"product_name"`
	ProductQuantity int    `json:"product_quantity"`
}
