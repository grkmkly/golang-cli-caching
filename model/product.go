package model

type Product struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Category             string    `json:"category"`
	Price                float64   `json:"price"`
	DiscountPercentage   float64   `json:"discountPercentage"`
	Rating               float64   `json:"rating"`
	Stock                int       `json:"stock"`
	Tags                 []string  `json:"tags"`
	Brand                string    `json:"brand,omitempty"`
	Sku                  string    `json:"sku"`
	Weight               int       `json:"weight"`
	Dimensions           Dimension `json:"dimensions"`
	WarrantyInformation  string    `json:"warrantyInformation"`
	ShippingInformation  string    `json:"shippingInformation"`
	AvailabilityStatus   string    `json:"availabilityStatus"`
	Reviews              []Review  `json:"reviews"`
	ReturnPolicy         string    `json:"returnPolicy"`
	MinimumOrderQuantity int       `json:"minimumOrderQuantity"`
	Metas                Meta      `json:"meta"`
	Images               []string  `json:"images"`
	Thumbnail            string    `json:"thumbnail"`
}

// https://mholt.github.io/json-to-go/ kullanılmıştır
