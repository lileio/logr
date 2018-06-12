package logr

import (
	"context"
	"fmt"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	zipkintracing "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"
)

const TraceKey = "traceId"

type Logr struct {
	ctx    context.Context
	Logger FieldLogr
}

type FieldLogr interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
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
	return &Logr{ctx: ctx, Logger: defaultLogger(ctx)}
}

func defaultLogger(ctx context.Context) FieldLogr {
	fields := logrus.Fields{}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		zs, ok := span.Context().(zipkintracing.SpanContext)
		if ok {
			traceID := zs.TraceID.ToHex()
			fields[TraceKey] = traceID
		}
	}

	lg := logrus.StandardLogger()
	lg.WithFields(fields)
	return lg
}

// LogToTrace logs the msg to the tracing span with the level
func (l *Logr) LogToTrace(level, msg string) {
	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		span.LogFields(log.String("event", msg))
	}
}

// LogErrorToTrace logs the msg to the tracing span with the level
func (l *Logr) LogErrorToTrace(level, msg string) {
	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		span.LogFields(log.String("stack", msg))
		span.SetTag("error", "true")
	}
}

func (l *Logr) Debugf(format string, args ...interface{}) {
	l.LogToTrace("DEBUG", fmt.Sprintf(format, args...))
	l.Logger.Debugf(format, args...)
}

func (l *Logr) Infof(format string, args ...interface{}) {
	l.LogToTrace("INFO", fmt.Sprintf(format, args...))
	l.Logger.Infof(format, args...)
}

func (l *Logr) Printf(format string, args ...interface{}) {
	l.LogToTrace("PRINT", fmt.Sprintf(format, args...))
	l.Logger.Printf(format, args...)
}

func (l *Logr) Warnf(format string, args ...interface{}) {
	l.LogToTrace("WARN", fmt.Sprintf(format, args...))
	l.Logger.Warnf(format, args...)
}

func (l *Logr) Warningf(format string, args ...interface{}) {
	l.LogToTrace("WARNING", fmt.Sprintf(format, args...))
	l.Logger.Warningf(format, args...)
}

func (l *Logr) Errorf(format string, args ...interface{}) {
	l.LogErrorToTrace("ERROR", fmt.Sprintf(format, args...))
	l.Logger.Errorf(format, args...)
}

func (l *Logr) Debug(args ...interface{}) {
	l.LogToTrace("DEBUG", fmt.Sprint(args...))
	l.Logger.Debug(args...)
}

func (l *Logr) DebugObject(name string, object interface{}) {
	span := opentracing.SpanFromContext(l.ctx)
	if span != nil {
		span.LogFields(log.Object(name, object))
	}

	l.Logger.Debug(name + ": " + litter.Sdump(object))
}

func (l *Logr) Info(args ...interface{}) {
	l.LogToTrace("INFO", fmt.Sprint(args...))
	l.Logger.Info(args...)
}

func (l *Logr) Print(args ...interface{}) {
	l.LogToTrace("PRINT", fmt.Sprint(args...))
	l.Logger.Print(args...)
}

func (l *Logr) Warn(args ...interface{}) {
	l.LogToTrace("WARN", fmt.Sprint(args...))
	l.Logger.Warn(args...)
}

func (l *Logr) Warning(args ...interface{}) {
	l.LogToTrace("WARNING", fmt.Sprint(args...))
	l.Logger.Warning(args...)
}

func (l *Logr) Error(args ...interface{}) {
	l.LogErrorToTrace("ERROR", fmt.Sprint(args...))
	l.Logger.Error(args...)
}
