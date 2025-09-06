package usecases

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/infra"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
)

type CreateOneLink struct {
	r              domain.LinkRepository
	processingChan chan uint
	wg             sync.WaitGroup
}

func NewCreateOneLinkUse(r domain.LinkRepository) *CreateOneLink {
	uc := &CreateOneLink{
		r:              r,
		processingChan: make(chan uint, 30),
	}

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		uc.wg.Add(1)
		go uc.worker(uint(i), uc.processingChan)
	}

	return uc
}

func (uc *CreateOneLink) Execute(r domain.NewLinkRequestDto) (domain.Link, error) {
	if err := r.Validate(); err != nil {
		return domain.Link{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeInvalidField,
			"CreateLink::Invalid::URL::%s",
			r.Url,
		)
	}

	link, err := uc.r.UpsertOne(r)

	if err != nil {
		return domain.Link{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"CreateLink::Create::Err::ValidateRequest",
		)
	}

	uc.processingChan <- link.ID
	logger.GetLogger().Debug("Link sent to processing channel", "linkID", link.ID)

	return link, nil
}

func (uc *CreateOneLink) Shutdown() {
	close(uc.processingChan)
	uc.wg.Wait()
	logger.GetLogger().Info("All workers have completed processing.")
}

type Summary struct {
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

func (uc *CreateOneLink) worker(id uint, jobs <-chan uint) {
	defer uc.wg.Done()

	logger.GetLogger().Debug("Worker started: ", "workerID", id)

	for linkId := range jobs {
		link, err := uc.r.GetOneById(domain.GetLinkRequestDto{ID: linkId})

		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("worker %d error fetching link %d: %v\n", id, linkId, err))
			continue
		}

		if link.Summary != "" {
			logger.GetLogger().Debug(fmt.Sprintf("worker %d skipping link %s as it already has a summary\n", id, link.Url))
			continue
		}

		// TODO: make categories a many to many relationship and store them
		// TODO: create a new use case to get links by category
		// TODO: create a new use case to get links by updatedAt desc
		result, err := infra.ExecSummarization(link.Url)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("worker %d error summarizing link %s: %v\n", id, link.Url, err))
			continue
		}

		logger.GetLogger().Debug(fmt.Sprintf("worker %d raw result for link %s: %s\n", id, link.Url, result))

		var summary Summary
		if err := json.Unmarshal([]byte(result), &summary); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("worker %d error unmarshaling result for link %s: %v\n", id, link.Url, err))
			continue
		}

		logger.GetLogger().Info(fmt.Sprintf("worker %d summary for link %s: %+v\n", id, link.Url, summary))

		if err := uc.r.UpdateSummary(link.ID, summary.Description); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("worker %d error updating summary for link %s: %v\n", id, link.Url, err))
			continue
		}
	}

}
