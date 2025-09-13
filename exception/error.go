package exception

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/pkg/errors"
	"runtime"
)

// PanicLogging logs an error with context and panics if the error is not nil
func PanicLogging(err interface{}) {
	if err != nil {
		// Get caller information for better context
		pc, file, line, _ := runtime.Caller(1)
		fn := runtime.FuncForPC(pc)

		// Convert interface to error if needed
		var errObj error
		switch e := err.(type) {
		case error:
			errObj = e
		case string:
			errObj = errors.New(e)
		default:
			errObj = errors.Errorf("%+v", err)
		}

		// Add stack trace if not already present
		if _, ok := errObj.(interface{ StackTrace() errors.StackTrace }); !ok {
			errObj = errors.WithStack(errObj)
		}

		// Log the error with the calling context
		logger.LogWithLevel(logger.Error, "Error detected", errObj,
			"caller_file", file,
			"caller_line", line,
			"caller_function", fn.Name())

		// Panic with the error
		panic(errObj)
	}
}
