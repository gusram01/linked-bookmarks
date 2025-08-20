package infra

import (
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/infra"
	"github.com/gusram01/linked-bookmarks/internal/link/usecases"
	"github.com/gusram01/linked-bookmarks/internal/platform/auth"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
	"github.com/gusram01/linked-bookmarks/internal/shared"
)

type LinkHandler struct {
	createOneUC  shared.QueryBaseUseCase[domain.NewLinkRequestDto, domain.Link]
	getOneByIdUC shared.QueryBaseUseCase[domain.GetLinkRequestDto, domain.Link]
	getAllUC     shared.QueryBaseUseCase[domain.GetAllLinksRequestDto, []domain.Link]
}

func newLinkHandler(
	createUc shared.QueryBaseUseCase[domain.NewLinkRequestDto, domain.Link],
	getOneUc shared.QueryBaseUseCase[domain.GetLinkRequestDto, domain.Link],
	getAll shared.QueryBaseUseCase[domain.GetAllLinksRequestDto, []domain.Link],
) *LinkHandler {
	return &LinkHandler{
		createOneUC:  createUc,
		getOneByIdUC: getOneUc,
		getAllUC:     getAll,
	}
}

func (lh *LinkHandler) registerRoutes(r fiber.Router) {
	lRouter := r.Group("api/links")

	lRouter.Use(auth.JwtClerkMiddleware())
	lRouter.Post("/", lh.createOne)
	lRouter.Get("/:id", lh.getOne)
	lRouter.Get("/", lh.getAll)
}

func Bootstrap(r fiber.Router) {
	repo := infra.NewLinkRepoWithGorm(database.DB)
	createUC := usecases.NewCreateOneLinkUse(repo)
	getOneUC := usecases.NewGetOneByIdLinkUse(repo)
	getAllUC := usecases.NewGetAllLinksUse(repo)

	newLinkHandler(createUC, getOneUC, getAllUC).registerRoutes(r)
}

func (lh *LinkHandler) createOne(c *fiber.Ctx) error {
	req := new(domain.NewLinkRequestDto)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
	}

	claims, err := auth.WithSessionClaims(c)
	if err != nil {
		return c.Status(401).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
	}

	req.Subject = claims.Subject

	link, ucErr := lh.createOneUC.Execute(*req)

	if ucErr != nil || link.ID == 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"error":   ucErr.Error(),
				"data":    nil,
			})
	}

	return c.Status(201).JSON(
		fiber.Map{
			"success": true,
			"error":   nil,
			"data": fiber.Map{
				"id":  link.ID,
				"url": link.Url,
			},
		},
	)

}

func (lh *LinkHandler) getOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid ID",
				"data":    nil,
			})
	}

	claims, cErr := auth.WithSessionClaims(c)

	if cErr != nil {
		return c.Status(401).JSON(
			fiber.Map{
				"success": false,
				"error":   cErr.Error(),
				"data":    nil,
			})
	}

	link, ucErr := lh.getOneByIdUC.Execute(domain.GetLinkRequestDto{
		ID:      uint(id),
		Subject: claims.Subject,
	})

	if ucErr != nil {
		return c.Status(404).JSON(
			fiber.Map{
				"success": false,
				"error":   ucErr.Error(),
				"data":    nil,
			})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"error":   nil,
			"data": fiber.Map{
				"id":  link.ID,
				"url": link.Url,
			},
		},
	)
}

func (lh *LinkHandler) getAll(c *fiber.Ctx) error {
	claims, err := auth.WithSessionClaims(c)
	if err != nil {
		return c.Status(401).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
	}

	loggedUser, uerr := user.Get(c.UserContext(), claims.Subject)

	// TODO: validate only can retrieve info from the current user
	if uerr != nil || loggedUser.ID != claims.Subject {
		return c.Status(403).JSON(
			fiber.Map{
				"success": false,
				"error":   "You do not have permission to access this resource",
				"data":    nil,
			})
	}

	var pagination domain.GetPaginatedLinksRequestDto

	if err := c.QueryParser(&pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			},
		)
	}

	req := new(domain.GetAllLinksRequestDto)

	req.Subject = claims.Subject
	req.Limit = pagination.PageSize
	req.Offset = pagination.PageNum * pagination.PageSize

	links, ucErr := lh.getAllUC.Execute(*req)

	if ucErr != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"success": false,
				"error":   ucErr.Error(),
				"data":    nil,
			})
	}

	var responseLinks []fiber.Map
	for _, link := range links {
		responseLinks = append(responseLinks, fiber.Map{
			"id":  link.ID,
			"url": link.Url,
		})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"error":   nil,
			"data":    responseLinks,
		},
	)
}
