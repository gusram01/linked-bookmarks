package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/ai"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"

	"github.com/iancoleman/strcase"
)

type SummarizeCategorizeLink struct {
	linkR domain.LinkRepository
	link  domain.Link
}

type Summary struct {
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

func NewSummarizeCategorizeLinkUse(r domain.LinkRepository, link domain.Link) *SummarizeCategorizeLink {
	return &SummarizeCategorizeLink{
		linkR: r,
		link:  link,
	}
}

func (uc *SummarizeCategorizeLink) Process() error {
	link, err := uc.linkR.GetOneById(domain.GetLinkRequestDto{ID: uc.link.ID})

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("worker error fetching link %d: %v\n", uc.link.ID, err))
		return internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::FetchLink::Err::%s",
			err.Error(),
		)
	}

	if link.Summary != "" {
		logger.GetLogger().Debug(fmt.Sprintf("worker skipping link %s as it already has a summary\n", link.Url))
		return internal.NewErrorf(
			internal.ErrorCodeSummaryExists,
			"SummarizeCategorizeLink::Skip::SummaryExists::LinkID::%d",
			link.ID,
		)
	}

	// TODO: create a new use case to get links by category
	// TODO: create a new use case to get links by updatedAt desc

	raw, rawErr := ai.MyAI.SummarizeAndCategorizeURL(link.Url)

	if rawErr != nil {
		return internal.WrapErrorf(
			rawErr,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::SummarizeAndCategorizeURL::Err::%s",
			rawErr.Error(),
		)
	}
	result, err := ai.MyAI.StructureOutput(raw)

	if err != nil {
		return internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::StructureOutput::Err::%s",
			err.Error(),
		)
	}

	logger.GetLogger().Debug(fmt.Sprintf("worker raw result for link %s: %s\n", link.Url, result))

	var summary Summary
	if err := json.Unmarshal([]byte(result), &summary); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("worker error unmarshaling result for link %s: %v\n", link.Url, err))
		return internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::Unmarshal::Err::%s",
			err.Error(),
		)
	}

	logger.GetLogger().Info(fmt.Sprintf("worker summary for link %s: %+v\n", link.Url, summary))

	if err := uc.linkR.UpdateSummary(domain.UpdateSummaryRequestDto{
		ID:      link.ID,
		Summary: summary.Description,
	}); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("worker error updating summary for link %s: %v\n", link.Url, err))
		return internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::UpdateSummary::Err::%s",
			err.Error(),
		)
	}

	if len(summary.Categories) == 0 {
		logger.GetLogger().Info(fmt.Sprintf("worker no categories found for link %s, skipping tag update\n", link.Url))
		return nil
	}

	snakeCaseTags := make([]string, len(summary.Categories))

	for i, tag := range summary.Categories {
		snakeCaseTags[i] = strcase.ToSnake(tag)
	}

	if err := uc.linkR.UpdateTags(domain.UpdateTagsRequestDto{
		ID:   link.ID,
		Tags: snakeCaseTags,
	}); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("worker error updating tags for link %s: %v\n", link.Url, err))
		return internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"SummarizeCategorizeLink::UpdateTags::Err::%s",
			err.Error(),
		)
	}

	return nil
}
