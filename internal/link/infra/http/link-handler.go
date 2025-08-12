package infra

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/infra"
	"github.com/gusram01/linked-bookmarks/internal/link/usecases"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
)


type LinkHandler struct {
    createOneUC usecases.BaseUseCase[domain.LinkRequest, domain.Link]
    getOneByIdUC usecases.BaseUseCase[uint, domain.Link]
}

func newLinkHandler(
    createUc usecases.BaseUseCase[domain.LinkRequest, domain.Link],
    getOneUc usecases.BaseUseCase[uint, domain.Link],
    ) *LinkHandler {
    return &LinkHandler{
        createOneUC: createUc,
        getOneByIdUC: getOneUc,
    }
}

func (lh *LinkHandler) registerRoutes(r fiber.Router) {
    lRouter := r.Group("api/links")

    lRouter.Post("/", lh.createOne)
    lRouter.Get("/:id", lh.getOne)

}

func Bootstrap(r fiber.Router) {
    repo := infra.NewLinkRepoWithGorm(database.DB)
    createUC := usecases.NewCreateOneLinkUse(repo)
    getOneUC := usecases.NewGetOneByIdLinkUse(repo)

    newLinkHandler(createUC, getOneUC).registerRoutes(r)
}


func (lh *LinkHandler) createOne(c *fiber.Ctx) error{
    req := new(domain.LinkRequest)

    if err := c.BodyParser(req); err != nil {
        return c.Status(400).JSON(
            fiber.Map{
                "success": false,
                "error": err.Error(),
                "data": nil,
            })
    }


    link, ucErr := lh.createOneUC.Execute(*req)

    if ucErr != nil {
        return c.Status(400).JSON(
            fiber.Map{
                "success": false,
                "error": ucErr.Error(),
                "data": nil,
            })
    }

    return c.Status(201).JSON(
        fiber.Map{
            "success": true,
            "error": nil,
            "data": fiber.Map{
                "id": link.ID,
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
                "error": "Invalid ID",
                "data": nil,
            })
    }

    link, ucErr := lh.getOneByIdUC.Execute(uint(id))

    if ucErr != nil {
        return c.Status(404).JSON(
            fiber.Map{
                "success": false,
                "error": ucErr.Error(),
                "data": nil,
            })
    }

    return c.Status(200).JSON(
        fiber.Map{
            "success": true,
            "error": nil,
            "data": fiber.Map{
                "id": link.ID,
                "url": link.Url,
            },
        },
    )
}

