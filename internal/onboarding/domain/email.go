package domain

import "time"

type Email struct {
	ID         uint      `json:"id"`
	ExternalID string    `json:"externalId"`
	Address    string    `json:"address"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  time.Time `json:"deletedAt"`
}
