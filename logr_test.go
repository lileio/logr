package logr_test

import (
	"context"
	"testing"

	"github.com/lileio/logr"
	"github.com/lileio/logr/logrfakes"
	opentracing "github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogr(t *testing.T) {
	l := logr.WithCtx(context.Background())
	assert.NotNil(t, l)
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

	fl := &logrfakes.FakeFieldLogger{}
	fl.WithFieldsReturns(logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	}))

	l := logr.WithCtx(ctx)
	l.Logger = fl
	l.Infof(":aasdasd")
	assert.NotNil(t, l)
}
