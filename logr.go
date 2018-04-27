package logr

import (
	"context"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	zipkintracing "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
)

const TraceKey = "traceId"

// SetLevelFromEnv is a convienience function to set the log level from
// env vars such as those on Kubernetes or Docker hosts
func SetLevelFromEnv() {
	lvl := os.Getenv("LOG_LEVEL")
	switch lvl {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// WithCtx will return a logrus logger that will log Zipkin when logging.
func WithCtx(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		zs, ok := span.Context().(zipkintracing.SpanContext)
		if ok {
			traceID := zs.TraceID.ToHex()
			fields[TraceKey] = traceID
		}
	}

	return logrus.WithFields(fields)
}

// ApplyCtx applies a traceid from a context an existing logrus logger
func ApplyCtx(ctx context.Context, e *logrus.Entry) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		zs, ok := span.Context().(zipkintracing.SpanContext)
		if ok {
			traceID := zs.TraceID.ToHex()
			e.Data[TraceKey] = traceID
		}
	}
}
