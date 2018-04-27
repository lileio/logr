package logr

import (
	"context"
	"fmt"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	zipkintracing "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
)

const TraceKey = "traceId"

type Logr struct {
	ctx    context.Context
	Logger logrus.FieldLogger
}

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
func WithCtx(ctx context.Context) *Logr {
	return &Logr{ctx: ctx}
}

// Logrus will return a logrus logger that will log the zipkin id when logging.
func (l *Logr) Logrus() *logrus.Entry {
	fields := logrus.Fields{}

	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		zs, ok := span.Context().(zipkintracing.SpanContext)
		if ok {
			traceID := zs.TraceID.ToHex()
			fields[TraceKey] = traceID
		}
	}

	if l.Logger == nil {
		l.Logger = logrus.StandardLogger()
	}

	return l.Logger.WithFields(fields)

}

// LogToTrace logs the msg to the tracing span with the level
func (l *Logr) LogToTrace(level, msg string) {
	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		span.LogFields(log.String("logr: "+level, msg))
	}
}

// LogErrorToTrace logs the msg to the tracing span with the level
func (l *Logr) LogErrorToTrace(level, msg string) {
	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		span.LogFields(log.String("logr: "+level, msg))
		span.SetTag("error", "true")
	}
}

func (l *Logr) Debugf(format string, args ...interface{}) {
	l.LogToTrace("DEBUG", fmt.Sprintf(format, args...))
	l.Logrus().Debugf(format, args...)
}

func (l *Logr) Infof(format string, args ...interface{}) {
	l.LogToTrace("INFO", fmt.Sprintf(format, args...))
	l.Logrus().Infof(format, args...)
}

func (l *Logr) Printf(format string, args ...interface{}) {
	l.LogToTrace("PRINT", fmt.Sprintf(format, args...))
	l.Logrus().Printf(format, args...)
}

func (l *Logr) Warnf(format string, args ...interface{}) {
	l.LogToTrace("WARN", fmt.Sprintf(format, args...))
	l.Logrus().Warnf(format, args...)
}

func (l *Logr) Warningf(format string, args ...interface{}) {
	l.LogToTrace("WARNING", fmt.Sprintf(format, args...))
	l.Logrus().Warningf(format, args...)
}

func (l *Logr) Errorf(format string, args ...interface{}) {
	l.LogErrorToTrace("ERROR", fmt.Sprintf(format, args...))
	l.Logrus().Errorf(format, args...)
}
