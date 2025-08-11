package greeter

import (
	"github.com/gofiber/fiber/v2"
)

type GrtService interface {
    Random(c *fiber.Ctx) error
    Predefined(c *fiber.Ctx) error
    PreWithResp(c *fiber.Ctx) error
}

type GrtServiceWithRepo struct {
    repo GrtRepo
}


func NewGreeterService(r GrtRepo) *GrtServiceWithRepo {
    return &GrtServiceWithRepo{
        repo: r,
    }
}

func (s *GrtServiceWithRepo) Predefined(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"ok": true, "greet": s.repo.Predefined()})
}

func (s *GrtServiceWithRepo) Random(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"ok": true, "greet": s.repo.Random()})
}

func (s *GrtServiceWithRepo) PreWithResp(c *fiber.Ctx) error {
    ans := c.Params("resp")

    return c.JSON(fiber.Map{"ok": true, "greet": s.repo.Answer(ans)})
}
