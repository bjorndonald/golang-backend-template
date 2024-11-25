package models

import "github.com/gofrs/uuid"

type GeoLocation struct {
	ID       uuid.UUID `json:"id,omitempty"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	IP       string    `json:"ip"`
	City     string    `json:"city"`
	Region   string    `json:"region"`
	Country  string    `json:"country"`
	Location string    `json:"loc"`
}
