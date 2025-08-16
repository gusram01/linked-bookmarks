package domain

import (
	"net/url"
	"time"
)

type UrlLink string

type Link struct {
	ID        uint
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type NewLinkRequestDto struct {
	Url     UrlLink `json:"url"`
	Subject string  `json:"subject"`
}

type GetLinkRequestDto struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
}

func (ul *NewLinkRequestDto) Validate() error {
	if _, err := url.ParseRequestURI(string(ul.Url)); err != nil {
		return err
	}

	return nil
}

type LinkResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type LinkRepository interface {
	Create(r NewLinkRequestDto) (Link, error)
	GetOneById(r GetLinkRequestDto) (Link, error)
	GetAll(cs string) ([]Link, error)
}
