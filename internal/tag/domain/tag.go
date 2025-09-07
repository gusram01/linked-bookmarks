package domain

import "time"

type Tag struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type CreateManyTagsRequestDto struct {
	Names  []string
	LinkID uint
}
