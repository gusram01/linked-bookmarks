package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/application/usecases"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/domain"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/infra"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
	"github.com/gusram01/linked-bookmarks/internal/shared"
)

type OnboardingHandler struct {
	upsertOneUserUC shared.CommandBaseUseCase[domain.NewUserRequestDto]
}

func newOnboardingHandler(upsertUc shared.CommandBaseUseCase[domain.NewUserRequestDto]) *OnboardingHandler {
	return &OnboardingHandler{
		upsertOneUserUC: upsertUc,
	}
}

func (uh *OnboardingHandler) registerRoutes(r fiber.Router) {
	oRouter := r.Group("api/onboarding")

	oRouter.Post("users/webhook", uh.upsertOneUser)
}

func Bootstrap(r fiber.Router) {
	repo := infra.NewUserRepoWithGorm(database.DB)
	upsertUc := usecases.NewUpsertUC(repo)
	newOnboardingHandler(upsertUc).registerRoutes(r)
}

func (uh *OnboardingHandler) upsertOneUser(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	payload := c.BodyRaw()

	user := new(domain.User)
	whEvent := new(domain.UserWebhookEvent)

	if err := c.BodyParser(whEvent); err != nil {
		logger.GetLogger().WarnContext(
			c.UserContext(),
			"Onboarding::Webhook::Invalid::Payload",
		)

		return c.SendStatus(400)
	}

	if err := uh.upsertOneUserUC.Execute(domain.NewUserRequestDto{
		User:      user,
		Event:     whEvent,
		RawHeader: headers,
		RawBody:   payload,
	}); err != nil {

		logger.GetLogger().Warn(
			err.Error(),
		)

		return c.SendStatus(400)
	}

	return c.SendStatus(fiber.StatusCreated)
}
