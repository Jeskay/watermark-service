package internal

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

func InitTracer(serviceName string, address string) (io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 100,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: address,
			LogSpans:           true,
		},
	}
	var closer io.Closer
	var err error
	Tracer, closer, err = cfg.NewTracer(config.Logger(jaeger.StdLogger))
	return closer, err
}

func StartSpan(name string, ctx context.Context) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span = Tracer.StartSpan(name, opentracing.ChildOf(span.Context()))
	} else {
		span = Tracer.StartSpan(name)
	}
	return span
}
