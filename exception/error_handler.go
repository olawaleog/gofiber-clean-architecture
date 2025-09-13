package exception

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// ErrorHandler is the central error handler for the fiber application
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Log error with stack trace
	logErrorWithStack(err)

	// Handle specific error types
	_, validationError := err.(ValidationError)
	if validationError {
		data := err.Error()
		var messages []map[string]interface{}

		errJson := json.Unmarshal([]byte(data), &messages)
		if errJson != nil {
			logger.LogWithLevel(logger.Error, "Failed to parse validation error JSON", errJson)
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Bad Request",
			Data:    messages,
		})
	}

	_, notFoundError := err.(NotFoundError)
	if notFoundError {
		return ctx.Status(fiber.StatusNotFound).JSON(model.GeneralResponse{
			Code:    404,
			Message: "Not Found",
			Data:    err.Error(),
		})
	}

	_, unauthorizedError := err.(UnauthorizedError)
	if unauthorizedError {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    err.Error(),
		})
	}

	_, badRequest := err.(BadRequestError)
	if badRequest {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: err.Error(),
			Data:    err.Error(),
		})
	}

	// If no specific error type matched, return a general error
	return ctx.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
		Code:    500,
		Message: err.Error(),
		Data:    "An unexpected error occurred",
	})
}

// logErrorWithStack logs the error with a stack trace
func logErrorWithStack(err error) {
	// Use pkg/errors.WithStack if not already wrapped
	if _, ok := err.(interface{ StackTrace() errors.StackTrace }); !ok {
		err = errors.WithStack(err)
	}

	// Log error with stack trace
	stackTrace := getStackTrace(err)
	logger.LogWithLevel(logger.Error, "Error occurred", err,
		"stack_trace", stackTrace)
}

// getStackTrace extracts and formats a stack trace from an error
func getStackTrace(err error) string {
	var stackTrace strings.Builder

	// Check if it's a pkg/errors error with stack trace
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err, ok := err.(stackTracer); ok {
		for i, f := range err.StackTrace() {
			// Only include up to 10 frames to avoid excessive logs
			if i >= 10 {
				break
			}

			pc := uintptr(f)
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				stackTrace.WriteString(fmt.Sprintf("%v: unknown\n", f))
				continue
			}

			file, line := fn.FileLine(pc)
			stackTrace.WriteString(fmt.Sprintf("%v: %s:%d - %s\n", f, file, line, fn.Name()))
		}
		return stackTrace.String()
	}

	// If it's not a pkg/errors, create a basic stack trace
	for i := 2; i < 12; i++ { // Skip first 2 frames (this function and caller)
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			stackTrace.WriteString(fmt.Sprintf("%s:%d - unknown\n", file, line))
			continue
		}
		stackTrace.WriteString(fmt.Sprintf("%s:%d - %s\n", file, line, fn.Name()))
	}

	return stackTrace.String()
}
