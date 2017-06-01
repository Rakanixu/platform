package handler

import (
	"errors"

	"cloud.google.com/go/translate"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"github.com/kazoup/platform/translate/srv/proto/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

const (
	quotaExceededMsg = "Quota for Translate service exceeded."
)

// Service ...
type Service struct{}

// Translate ...
func (s *Service) Translate(
	ctx context.Context,
	req *proto_translate.TranslateRequest,
	res *proto_translate.TranslateResponse) error {

	if err := validate.Exists(
		req.GetSourceLang(),
		req.GetDestLang(),
	); err != nil {
		return err
	}

	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	_, _, rate, _, quota, ok := quota.Check(ctx, globals.TRANSLATE_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			"",
			403,
			errors.New("quota.Check failed"),
		)
	}

	if rate-quota > 0 {
		res.Info = quotaExceededMsg
		return nil
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			"translate.NewClient failed",
			403,
			err,
		)
	}

	src, err := language.Parse(req.GetSourceLang())
	if err != nil {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			`language.Parse("`+req.GetSourceLang()+`") failed`,
			403,
			err,
		)
	}

	dst, err := language.Parse(req.GetDestLang())
	if err != nil {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			`language.Parse("`+req.GetDestLang()+`") failed`,
			403,
			err,
		)
	}

	translations, err := client.Translate(
		ctx,
		req.GetText(),
		dst,
		&translate.Options{
			Source: src,
			Format: translate.Text,
		},
	)

	if err != nil {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			"client.Translate failed",
			403,
			err,
		)
	}

	for _, val := range translations {
		res.Translations = append(res.Translations, val.Text)
	}

	if err := client.Close(); err != nil {
		return platform_errors.NewPlatformError(
			globals.TRANSLATE_SERVICE_NAME,
			"Translate",
			"client.Close failed",
			403,
			err,
		)
	}

	return nil
}
