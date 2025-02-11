package middlewares

import (
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Recover(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			funcName := extractFuncNameBeforePanic()
			log.Error().Str("funcName", funcName).Interface("recover", r).
				Str("stackTrace", string(debug.Stack())).
				Interface("request", string(c.Request().Body())).
				Msg("panic occurred")
			c.Context().Error("Internal Server Error", http.StatusInternalServerError)
		}
	}()
	return c.Next()
}

func extractFuncNameBeforePanic() string {
	var funcName, funcNamePrev string
	for i := 1; ; i++ {
		if pc, _, _, ok := runtime.Caller(i); !ok {
			break
		} else if details := runtime.FuncForPC(pc); details == nil {
			break
		} else {
			funcNamePrev, funcName = funcName, details.Name()
		}
		if strings.HasSuffix(funcNamePrev, ".gopanic") {
			break
		}
	}
	return funcName
}
