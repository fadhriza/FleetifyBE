package middleware

import (
	"fmt"

	"fleetify/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			errors.LogError("Panic Recovered", fmt.Errorf("%v", e))
		},
	})
}
