package usecases

import (
	"context"
	"strconv"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
	vectordb "github.com/gusram01/linked-bookmarks/internal/vector-db"
)

type SemanticLinksSearch struct {
	r domain.LinkRepository
}

func NewSemanticLinksSearchUse(r domain.LinkRepository) *SemanticLinksSearch {
	return &SemanticLinksSearch{
		r: r,
	}
}

type VectorResponse struct {
	IDs       []string                 `json:"ids"`
	Documents []string                 `json:"documents"`
	Metadatas []map[string]interface{} `json:"metadatas"`
	Distances []float64                `json:"distances"`
	Include   []string                 `json:"include"`
}

func (uc *SemanticLinksSearch) Execute(req domain.SemanticSearchRequestDto) ([]domain.Link, error) {
	ctx := context.Background()

	res, err := vectordb.VDB.Collection.Query(
		ctx,
		chroma.WithQueryTexts(req.Query),
		chroma.WithNResults(5),
	)

	if err != nil {
		return []domain.Link{}, err
	}

	logger.GetLogger().Info("Semantic search results: ", "response", res)

	ids := res.GetIDGroups()

	if len(ids) == 0 || len(ids[0]) == 0 {
		return []domain.Link{}, nil
	}

	var ucReq domain.GetManyLinksByIdsRequestDto

	ucReq.IDs = make([]uint, len(ids[0]))
	for i, id := range ids[0] {
		parsedID, err := strconv.ParseUint(string(id), 10, 64)
		if err != nil {
			logger.GetLogger().Error("Error parsing ID: ", "id", id, "error", err)

			continue
		}
		ucReq.IDs[i] = uint(parsedID)
	}

	ucReq.Subject = req.Subject

	links, err := uc.r.GetManyByIds(ucReq)
	if err != nil {
		return []domain.Link{}, err
	}

	return links, nil

}
