package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get("X-Request-ID")

			if reqID == "" {
				reqID = uuid.New().String()
			}

			c.Response().Header().Set("X-Request-ID", reqID)
			c.Set("request_id", reqID)
			return next(c)
		}
	}
}

func LoggingMiddleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			req := c.Request()
			res := c.Response()

			requestID, _ := c.Get("request_id").(string)

			baseFields := []interface{}{
				"request_id", requestID,
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"remote_ip", c.RealIP(),
				"user_agent", req.UserAgent(),
				"latency", duration.String(),
				"latency_ms", duration.Milliseconds(),
				"bytes_out", res.Size,
			}

			if err != nil {
				logger.Errorw("Request processing error",
					append(baseFields, "error", err.Error())...,
				)
			} else {
				logger.Infow("Request completed", baseFields...)
			}

			return err
		}
	}
}
