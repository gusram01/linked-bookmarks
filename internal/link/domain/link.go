package domain

import (
	"net/url"
	"time"
)

type UrlLink string

type Link struct {
    ID uint
    Url string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time
}

type LinkRequest struct {
    Url UrlLink `json:"url"`;
}

func (ul *LinkRequest) Validate() error {
    if _, err := url.ParseRequestURI(string(ul.Url)); err != nil {
        return err
    }

    return nil
}

type LinkResponse struct {
    Success bool `json:"success"`
    Data interface{} `json:"data,omitempty"`
    Error interface{} `json:"error,omitempty"`
}

type LinkRepository interface {
    Create(r LinkRequest) (Link, error)
    GetOneById(id uint) (Link, error)
    GetAll() ([]Link, error)
}
