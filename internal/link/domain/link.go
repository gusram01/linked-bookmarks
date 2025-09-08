package domain

import (
	"net/url"
	"time"

	"github.com/gusram01/linked-bookmarks/internal"
)

type UrlLink string

type Link struct {
	ID        uint      `json:"id"`
	Url       string    `json:"url"`
	Summary   string    `json:"summary"`
	Attempts  uint      `json:"attempts"`
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

	if ul.Subject == "" {
		return internal.NewErrorf(
			internal.ErrorCodeInvalidField,
			"LinkRequest::Subject::Invalid",
		)
	}

	return nil
}

type GetLinkRequestDto struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
}

type GetAllLinksRequestDto struct {
	Subject string `json:"subject"`
	Limit   uint   `json:"limit"`
	Offset  uint   `json:"offset"`
}

type GetAllQueryResultDto struct {
	Links      []Link `json:"links"`
	TotalCount int64  `json:"totalCount"`
	Pages      int64  `json:"pages"`
}

type GetAllLinksResponseDto struct {
	Links      []any `json:"links"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int64 `json:"totalPages"`
}

type GetPaginatedLinksRequestDto struct {
	PageNum  uint `query:"pageNum"`
	PageSize uint `query:"pageSize"`
}

type UpdateSummaryRequestDto struct {
	ID      uint   `json:"id"`
	Summary string `json:"summary"`
}

type UpdateTagsRequestDto struct {
	ID   uint     `json:"id"`
	Tags []string `json:"tags"`
}
