package infra

import (
	"errors"

	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/infra"
	"github.com/gusram01/linked-bookmarks/internal/link/usecases"
	"github.com/gusram01/linked-bookmarks/internal/platform/auth"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
	"github.com/gusram01/linked-bookmarks/internal/shared"
	"github.com/gusram01/linked-bookmarks/internal/worker"
)

type LinkHandler struct {
	createOneUC      shared.QueryBaseUseCase[domain.NewLinkRequestDto, domain.Link]
	getOneByIdUC     shared.QueryBaseUseCase[domain.GetLinkRequestDto, domain.Link]
	getAllUC         shared.QueryBaseUseCase[domain.GetAllLinksRequestDto, domain.GetAllQueryResultDto]
	semanticSearchUC shared.QueryBaseUseCase[domain.SemanticSearchRequestDto, []domain.Link]
	linkRepo         domain.LinkRepository
}

func newLinkHandler(
	createUc shared.QueryBaseUseCase[domain.NewLinkRequestDto, domain.Link],
	getOneUc shared.QueryBaseUseCase[domain.GetLinkRequestDto, domain.Link],
	getAll shared.QueryBaseUseCase[domain.GetAllLinksRequestDto, domain.GetAllQueryResultDto],
	semanticSearchUc shared.QueryBaseUseCase[domain.SemanticSearchRequestDto, []domain.Link],
	linkRepo domain.LinkRepository,

) *LinkHandler {
	return &LinkHandler{
		createOneUC:      createUc,
		getOneByIdUC:     getOneUc,
		getAllUC:         getAll,
		semanticSearchUC: semanticSearchUc,
		linkRepo:         linkRepo,
	}
}

func (lh *LinkHandler) registerRoutes(r fiber.Router) {
	lRouter := r.Group("api/links")

	lRouter.Use(auth.JwtClerkMiddleware())
	lRouter.Post("/", lh.createOne)
	lRouter.Get("/search", lh.searchLinks)
	lRouter.Get("/search/:id", lh.getOne)
	lRouter.Get("/", lh.getAll)
}

func Bootstrap(r fiber.Router) {
	repo := infra.NewLinkRepoWithGorm(database.DB)
	createUC := usecases.NewCreateOneLinkUse(repo)
	getOneUC := usecases.NewGetOneByIdLinkUse(repo)
	getAllUC := usecases.NewGetAllLinksUse(repo)
	semanticSearchUC := usecases.NewSemanticLinksSearchUse(repo)

	newLinkHandler(createUC, getOneUC, getAllUC, semanticSearchUC, repo).registerRoutes(r)
}

func (lh *LinkHandler) createOne(c *fiber.Ctx) error {
	req := new(domain.NewLinkRequestDto)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(internal.NewGcResponse(nil, err))
	}

	claims, err := auth.WithSessionClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(internal.NewGcResponse(nil, err))
	}

	req.Subject = claims.Subject

	link, ucErr := lh.createOneUC.Execute(*req)

	if ucErr != nil || link.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(internal.NewGcResponse(nil, ucErr))
	}

	scUC := usecases.NewSummarizeCategorizeLinkUse(lh.linkRepo, link)

	worker.CentralWorkerPool.Submit(scUC)

	return c.Status(fiber.StatusCreated).JSON(
		internal.NewGcResponse(
			internal.GcMap{
				"id":  link.ID,
				"url": link.Url,
			},
			nil,
		),
	)

}

func (lh *LinkHandler) searchLinks(c *fiber.Ctx) error {
	claims, err := auth.WithSessionClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(internal.NewGcResponse(nil, err))
	}

	query := c.Query("s", "")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(internal.NewGcResponse(nil, errors.New("query parameter 's' is required")))
	}

	linksResult, err := lh.semanticSearchUC.Execute(domain.SemanticSearchRequestDto{
		Query:   query,
		Subject: claims.Subject,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internal.NewGcResponse(nil, err))
	}

	var responseLinks domain.GetAllLinksResponseDto

	responseLinks.TotalCount = int64(len(linksResult))
	responseLinks.TotalPages = 1

	for _, link := range linksResult {
		responseLinks.Links = append(
			responseLinks.Links,
			internal.GcMap{
				"id":       link.ID,
				"url":      link.Url,
				"summary":  link.Summary,
				"attempts": link.Attempts,
			})
	}

	return c.Status(fiber.StatusOK).JSON(internal.NewGcResponse(responseLinks, nil))

}

func (lh *LinkHandler) getOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			internal.NewGcResponse(
				nil,
				errors.New("invalid id"),
			))
	}

	claims, cErr := auth.WithSessionClaims(c)

	if cErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(internal.NewGcResponse(nil, cErr))
	}

	link, ucErr := lh.getOneByIdUC.Execute(domain.GetLinkRequestDto{
		ID:      uint(id),
		Subject: claims.Subject,
	})

	if ucErr != nil {
		return c.Status(fiber.StatusNotFound).JSON(internal.NewGcResponse(nil, ucErr))
	}

	return c.Status(fiber.StatusOK).JSON(
		internal.NewGcResponse(
			internal.GcMap{
				"id":       link.ID,
				"url":      link.Url,
				"summary":  link.Summary,
				"attempts": link.Attempts,
			},
			nil,
		),
	)
}

func (lh *LinkHandler) getAll(c *fiber.Ctx) error {
	claims, err := auth.WithSessionClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(internal.NewGcResponse(nil, err))
	}

	loggedUser, uerr := user.Get(c.UserContext(), claims.Subject)

	// TODO: validate only can retrieve info from the current user
	if uerr != nil || loggedUser.ID != claims.Subject {
		return c.Status(fiber.StatusForbidden).JSON(internal.NewGcResponse(
			nil,
			errors.New("you do not have permission to access this resource"),
		))
	}

	pagination := new(domain.GetPaginatedLinksRequestDto)

	if err := c.QueryParser(pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(internal.NewGcResponse(nil, err))
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 5
	}

	req := new(domain.GetAllLinksRequestDto)

	req.Subject = claims.Subject
	req.Limit = pagination.PageSize
	req.Offset = pagination.PageNum * pagination.PageSize

	linksResult, ucErr := lh.getAllUC.Execute(*req)

	if ucErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internal.NewGcResponse(nil, ucErr))
	}

	var responseLinks domain.GetAllLinksResponseDto

	responseLinks.TotalCount = linksResult.TotalCount
	responseLinks.TotalPages = linksResult.Pages

	for _, link := range linksResult.Links {
		responseLinks.Links = append(
			responseLinks.Links,
			internal.GcMap{
				"id":       link.ID,
				"url":      link.Url,
				"summary":  link.Summary,
				"attempts": link.Attempts,
			})
	}

	return c.Status(fiber.StatusOK).JSON(internal.NewGcResponse(responseLinks, nil))
}
