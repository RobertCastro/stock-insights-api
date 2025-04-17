package models

import (
	"time"
)

// Representa la información de una acción con sus ratings
type Stock struct {
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}
