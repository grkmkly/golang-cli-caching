package model

import "time"

type Meta struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Barcode   string    `json:"barcode"`
	QrCode    string    `json:"qrCode"`
}
