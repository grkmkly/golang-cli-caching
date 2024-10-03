package model

import "time"

type Review struct {
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment"`
	Date          time.Time `json:"date"`
	ReviewerName  string    `json:"reviewerName"`
	ReviewerEmail string    `json:"reviewerEmail"`
}
