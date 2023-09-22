package fibermid

import (
	"github.com/gofiber/fiber/v2"

	"github.com/sandwich-go/xmonitor"
)

// SkipFunc if skip monitor return true else return false
type SkipFunc func(*fiber.Ctx) bool

func NewMonitorMid(skipper SkipFunc, collector xmonitor.Collector) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if skipper != nil && skipper(ctx) {
			return ctx.Next()
		}
		method := ctx.Route().Method
		path := ctx.Route().Path
		after := collector.MonitorRequest(method, path, 0)
		err := ctx.Next()
		// initialize with default error code
		// https://docs.gofiber.io/guide/error-handling
		status := fiber.StatusInternalServerError
		if err != nil {
			if e, ok := err.(*fiber.Error); ok {
				// Get correct error code from fiber.Error type
				status = e.Code
			}
		} else {
			status = ctx.Response().StatusCode()
		}
		after(status, 0)
		return err
	}
}
