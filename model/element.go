package model

type Element struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Skip     int       `json:"skip"`
	Limit    int       `json:"limit"`
}
