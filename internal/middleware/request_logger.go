package middleware

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/pkg/utils"
)

const requestIDLocalsKey = "requestid"

var skipBodyPaths = map[string]struct{}{
	"/health":       {},
	"/docs":         {},
	"/swagger.json": {},
}

var sampledPaths = map[string]struct{}{
	"/health": {},
}

var (
	healthLogMu       sync.Mutex
	healthLogCounters = make(map[string]*uint64)
)

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if rid, ok := c.Locals(requestIDLocalsKey).(string); ok && rid != "" {
			c.SetUserContext(utils.WithRequestID(c.UserContext(), rid))
			c.Context().SetUserValue(utils.RequestIDKey, rid)
		}
		return c.Next()
	}
}

func AccessLog(logBody bool, bodyMaxBytes, healthSampleN int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		captureBody := logBody && shouldLogBody(c.Path())
		var reqBody any
		if captureBody {
			reqBody = bodyForLog(c.Get(fiber.HeaderContentType), nil, c.Body(), bodyMaxBytes)
		}

		err := c.Next()
		status := c.Response().StatusCode()
		if err != nil && status < fiber.StatusBadRequest {
			status = fiber.StatusInternalServerError
		}
		if shouldSampleOut(c.Path(), status, healthSampleN, c.IP()) {
			return err
		}

		latencyMs := float64(time.Since(start).Microseconds()) / 1000
		args := []any{
			"request_id", c.Locals(requestIDLocalsKey),
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency_ms", latencyMs,
			"ip", c.IP(),
			"bytes", len(c.Response().Body()),
		}
		if uid := c.Locals("user_id"); uid != nil {
			args = append(args, "user_id", uid)
		}
		if q := string(c.Request().URI().QueryString()); q != "" {
			args = append(args, "query", q)
		}
		if captureBody {
			if reqBody != nil {
				args = append(args, "request_body", reqBody)
			}
			respBody := bodyForLog(string(c.Response().Header.ContentType()), c.Response().Header.Peek(fiber.HeaderContentEncoding), c.Response().Body(), bodyMaxBytes)
			if respBody != nil {
				args = append(args, "response_body", respBody)
			}
		}

		utils.Log("Access").Info(fmt.Sprintf("%s %s (%.3fms)", c.Method(), c.Path(), latencyMs), args...)
		return err
	}
}

func shouldLogBody(path string) bool {
	_, skip := skipBodyPaths[path]
	return !skip
}

func shouldSampleOut(path string, status, sampleN int, ip string) bool {
	if sampleN <= 1 || status >= fiber.StatusBadRequest {
		return false
	}
	if _, ok := sampledPaths[path]; !ok {
		return false
	}
	healthLogMu.Lock()
	ctr, ok := healthLogCounters[ip]
	if !ok {
		ctr = new(uint64)
		healthLogCounters[ip] = ctr
	}
	healthLogMu.Unlock()
	n := atomic.AddUint64(ctr, 1)
	return (n-1)%uint64(sampleN) != 0
}

func bodyForLog(contentType string, contentEncoding, body []byte, maxBytes int) any {
	if len(body) == 0 {
		return nil
	}
	if len(contentEncoding) > 0 {
		return "[" + string(contentEncoding) + " omitted]"
	}
	if !strings.Contains(strings.ToLower(contentType), "application/json") {
		return "[" + mediaType(contentType) + " omitted]"
	}
	return utils.RedactJSONValue(body, maxBytes)
}

func mediaType(contentType string) string {
	mt := contentType
	if i := strings.IndexByte(mt, ';'); i >= 0 {
		mt = mt[:i]
	}
	if mt = strings.TrimSpace(mt); mt == "" {
		return "binary"
	}
	return mt
}
