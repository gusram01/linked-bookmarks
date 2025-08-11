package greeter

import (
	"github.com/gofiber/fiber/v2"
)

type GreeterHandler struct {
    svc GrtService
}

func newGreeterHandler(s GrtService) *GreeterHandler {
    return &GreeterHandler{
        svc: s,
    }
}

func (s *GreeterHandler) registerRoutes(r fiber.Router) {
    greets := r.Group("/api/greet")
    greets.Get("/", s.svc.Predefined)
    greets.Get("/random", s.svc.Random)
    greets.Get("/interact/:resp?", s.svc.PreWithResp)
}

func Bootstrap(c fiber.Router){
    repo := NewGrtRepo()
    svc := NewGreeterService(repo)

    newGreeterHandler(svc).registerRoutes(c)
}
