package domain

import (
	"net/url"
	"time"
)

type UrlLink string

type Link struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type NewLinkRequestDto struct {
	Url     UrlLink `json:"url"`
	Subject string  `json:"subject"`
}

func (ul *NewLinkRequestDto) Validate() error {
	if _, err := url.ParseRequestURI(string(ul.Url)); err != nil {
		return err
	}

	return nil
}

type GetLinkRequestDto struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
}
