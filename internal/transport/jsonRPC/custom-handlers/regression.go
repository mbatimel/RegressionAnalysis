package custom_handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	customErrors "github.com/mbatimel/RegressionAnalysis/internal/errors"
	"github.com/mbatimel/RegressionAnalysis/pkg/models/consts"
)

type RestResponse struct {
	Data             interface{}            `json:"data"`
	Error            bool                   `json:"error"`
	ErrorText        string                 `json:"errorText"`
	AdditionalErrors map[string]interface{} `json:"additionalErrors"`
}

func sendResponse(ctx *fiber.Ctx, log zerolog.Logger, data interface{}, respError error) {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(http.StatusOK)

	response := &RestResponse{
		Data:  data,
		Error: respError != nil,
	}

	if response.Error {
		customErr, ok := respError.(*customErrors.Error)
		if !ok {
			switch {
			case strings.Contains(respError.Error(), customErrors.AccessDeniedError().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusForbidden)
				response.ErrorText = consts.ErrForbidden
			case strings.Contains(respError.Error(), customErrors.BadMetaValue().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusBadRequest)
				response.ErrorText = consts.ErrBadRequest
			default:
				ctx.Response().SetStatusCode(fasthttp.StatusInternalServerError)
				response.ErrorText = consts.ErrInternal
			}
		} else {
			switch {
			case strings.Contains(customErr.ErrorText, customErrors.InternalServerError().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusInternalServerError)
				response.ErrorText = customErr.ErrorText
			case strings.Contains(customErr.ErrorText, customErrors.MethodNotAllowedError().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusBadRequest)
				response.ErrorText = customErr.ErrorText
			case strings.Contains(customErr.ErrorText, customErrors.MethodNotAllowedError().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusMethodNotAllowed)
				response.ErrorText = customErr.ErrorText
			case strings.Contains(customErr.ErrorText, customErrors.ForbiddenError().Error()):
				ctx.Response().SetStatusCode(fasthttp.StatusForbidden)
				response.ErrorText = customErr.ErrorText
			default:
				ctx.Response().SetStatusCode(fasthttp.StatusInternalServerError)
				response.ErrorText = consts.ErrInternal

			}

			if customErr.GetOuterError() != nil {
				switch {
				case strings.Contains(customErr.GetOuterError().Error(), customErrors.AccessDeniedError().ErrorText):
					ctx.Response().SetStatusCode(fasthttp.StatusForbidden)
					response.ErrorText = consts.ErrAccessDenied
				case strings.Contains(customErr.GetOuterError().Error(), customErrors.MethodNotAllowedError().ErrorText):
					ctx.Response().SetStatusCode(fasthttp.StatusBadRequest)
					response.ErrorText = consts.ErrMethodNotAllowed
				case strings.Contains(customErr.GetOuterError().Error(), customErrors.AlreadyExists().ErrorText):
					ctx.Response().SetStatusCode(fasthttp.StatusConflict)
					response.ErrorText = consts.ErrBadRequest
				}
			}
		}

		if customErr != nil && customErr.Cause != nil {
			response.AdditionalErrors = customErr.Cause
		}
	}

	respBody, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return
	}

	if _, err = ctx.Write(respBody); err != nil {
		log.Error().Err(err).Msg("failed to send response")
		return

	}
}
