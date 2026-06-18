package utils

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

var (
	baseHandler   slog.Handler
	extraHandlers []slog.Handler
	quiet         = true
	logLevel      = new(slog.LevelVar)
)

var quietAllowModules = map[string]struct{}{
	"Access": {},
	"App":    {},
}

func InitLogger() {
	_ = os.MkdirAll("logs", 0755)
	writer, err := NewDailyRotateWriter("logs/app")
	if err != nil {
		panic(err)
	}
	baseHandler = slog.NewJSONHandler(io.MultiWriter(writer, os.Stdout), &slog.HandlerOptions{Level: logLevel})
	rebuildLogger()
}

func AttachHandler(h slog.Handler) {
	if h == nil {
		return
	}
	extraHandlers = append(extraHandlers, h)
	rebuildLogger()
}

func SetQuiet(q bool) {
	quiet = q
	rebuildLogger()
}

func rebuildLogger() {
	if baseHandler == nil {
		return
	}
	var h slog.Handler = baseHandler
	if len(extraHandlers) > 0 {
		h = fanoutHandler{handlers: append([]slog.Handler{baseHandler}, extraHandlers...)}
	}
	if quiet {
		h = noiseFilterHandler{inner: h}
	}
	Logger = slog.New(h)
}

type noiseFilterHandler struct {
	inner slog.Handler
}

func (h noiseFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h noiseFilterHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.inner.Handle(ctx, r)
	}
	allowed := false
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "module" {
			_, allowed = quietAllowModules[a.Value.String()]
			return false
		}
		return true
	})
	if allowed {
		return h.inner.Handle(ctx, r)
	}
	return nil
}

func (h noiseFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return noiseFilterHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h noiseFilterHandler) WithGroup(name string) slog.Handler {
	return noiseFilterHandler{inner: h.inner.WithGroup(name)}
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func SetLogLevel(level string) {
	logLevel.Set(ParseLevel(level))
}

type fanoutHandler struct {
	handlers []slog.Handler
}

func (f fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range f.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (f fanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range f.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		next[i] = h.WithAttrs(attrs)
	}
	return fanoutHandler{handlers: next}
}

func (f fanoutHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		next[i] = h.WithGroup(name)
	}
	return fanoutHandler{handlers: next}
}

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

type ModuleLogger struct {
	module  string
	ctxArgs []any
}

func Log(module string) *ModuleLogger {
	return &ModuleLogger{module: module}
}

func LogCtx(ctx context.Context, module string) *ModuleLogger {
	l := &ModuleLogger{module: module}
	if rid := RequestIDFromContext(ctx); rid != "" {
		l.ctxArgs = []any{"request_id", rid}
	}
	return l
}

func (l *ModuleLogger) prepend(args []any) []any {
	out := make([]any, 0, len(l.ctxArgs)+len(args)+2)
	out = append(out, "module", l.module)
	out = append(out, l.ctxArgs...)
	out = append(out, args...)
	return out
}

func (l *ModuleLogger) Debug(msg string, args ...any) {
	Logger.Debug(msg, l.prepend(args)...)
}

func (l *ModuleLogger) Info(msg string, args ...any) {
	Logger.Info(msg, l.prepend(args)...)
}

func (l *ModuleLogger) Warn(msg string, args ...any) {
	Logger.Warn(msg, l.prepend(args)...)
}

func (l *ModuleLogger) Error(msg string, args ...any) {
	for i := 0; i+1 < len(args); i += 2 {
		if key, ok := args[i].(string); ok && key == "error" {
			if errVal, ok := args[i+1].(error); ok {
				args[i+1] = errVal.Error()
			}
		}
	}
	Logger.Error(msg, l.prepend(args)...)
}
