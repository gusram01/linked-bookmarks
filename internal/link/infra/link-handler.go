package infra

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/usecases"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
)


type LinkHandler struct {
    createOneUC usecases.CreateOneLink
}

func newLinkHandler(uc usecases.CreateOneLink) *LinkHandler {
    return &LinkHandler{
        createOneUC: uc,
    }
}

func (lh *LinkHandler) registerRoutes(r fiber.Router) {
    lRouter := r.Group("api/links")

    lRouter.Post("/", lh.createOne)

}

func Bootstrap(r fiber.Router) {

    repo := NewLinkRepoWithGorm(database.DB)
    createUC := usecases.NewCreateOneLinkUse(repo)

    newLinkHandler(*createUC).registerRoutes(r)
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
            "data": link.ID,
        },
    )

}

