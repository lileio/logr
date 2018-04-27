package logr_test

import (
	"context"
	"testing"

	"github.com/lileio/logr"
	opentracing "github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogr(t *testing.T) {
	l := logr.WithCtx(context.Background())
	assert.NotNil(t, l)
	assert.Len(t, l.Data, 0)
}

func TestLogrWithZipkinContext(t *testing.T) {
	recorder := zipkintracer.NewInMemoryRecorder()
	tracer, err := zipkintracer.NewTracer(recorder)
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("test_logger")
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	if err != nil {
		t.Fatal(err)
	}

	l := logr.WithCtx(ctx)
	assert.NotNil(t, l)
	assert.Len(t, l.Data, 1)
}

func TestApply(t *testing.T) {
	recorder := zipkintracer.NewInMemoryRecorder()
	tracer, err := zipkintracer.NewTracer(recorder)
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("test_logger")
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	if err != nil {
		t.Fatal(err)
	}

	l := logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	})

	logr.ApplyCtx(ctx, l)

	assert.NotNil(t, l)
	assert.Len(t, l.Data, 2)
}
