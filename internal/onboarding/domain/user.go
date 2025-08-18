package domain

import "time"

type User struct {
	ID        uint      `json:"id"`
	AuthID    string    `json:"authId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type NewUserRequestDto struct {
	User      *User
	Event     *UserWebhookEvent
	RawHeader map[string][]string
	RawBody   []byte
}
