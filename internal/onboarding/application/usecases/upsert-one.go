package usecases

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/domain"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	svix "github.com/svix/svix-webhooks/go"
)

var secret = config.ENVS.WebhookProviderSecret

type UpsertOneUserUC struct {
	r  domain.UserRepository
	wh *svix.Webhook
}

type BaseData struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

func NewUpsertUC(repo domain.UserRepository) *UpsertOneUserUC {

	wh, err := svix.NewWebhook(secret)

	if err != nil {
		log.Fatal(err)
	}

	return &UpsertOneUserUC{
		r:  repo,
		wh: wh,
	}
}

func (uc *UpsertOneUserUC) Execute(req domain.NewUserRequestDto) error {

	if err := uc.wh.Verify(req.RawBody, req.RawHeader); err != nil {
		return internal.WrapErrorf(
			err,
			internal.ErrorCodePermissionDenied,
			"Webhook::Validation::failed",
		)
	}

	var baseData BaseData

	if err := json.Unmarshal(req.Event.Data, &baseData); err != nil {
		return internal.NewErrorf(
			internal.ErrorCodeDataLoss,
			"Webhook::Data::Invalid",
		)
	}

	if baseData.ID == "" {
		return internal.NewErrorf(
			internal.ErrorCodeDataLoss,
			"Webhook::Id::Lost",
		)
	}

	req.User.AuthID = baseData.ID

	if err := uc.r.Upsert(req.User); err != nil {
		var ierr *internal.Error

		if !errors.As(err, &ierr) {
			return err
		}

		switch ierr.Code() {
		case internal.ErrorCodeWHHandleUserFound:
			return nil
		default:
			return ierr.Unwrap()
		}
	}

	return nil
}
