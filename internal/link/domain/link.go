package domain

import "net/url"

type UrlLink string

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
    Create(r LinkRequest) error
    GetOneById(id string) (LinkResponse, error)
    GetAll() (LinkResponse, error)
}
