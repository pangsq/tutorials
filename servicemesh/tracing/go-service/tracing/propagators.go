package tracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type b3Propagator struct {
	b3.HTTPFormat
}

func (propagator b3Propagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}
	carrier.Set(b3.TraceIDHeader, sc.TraceID().String())
	carrier.Set(b3.SpanIDHeader, sc.SpanID().String())

	var sampled string
	if sc.IsSampled() {
		sampled = "1"
	} else {
		sampled = "0"
	}
	carrier.Set(b3.SampledHeader, sampled)
}

func (propagator b3Propagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	tid, ok := b3.ParseTraceID(carrier.Get(b3.TraceIDHeader))
	if !ok {
		return ctx
	}
	sid, ok := b3.ParseSpanID(carrier.Get(b3.SpanIDHeader))
	if !ok {
		return ctx
	}
	sampled, _ := b3.ParseSampled(carrier.Get(b3.SampledHeader))
	return trace.ContextWithSpanContext(ctx, trace.SpanContext{}.
		WithTraceID(trace.TraceID(tid)).
		WithSpanID(trace.SpanID(sid)).
		WithTraceFlags(trace.TraceFlags(sampled)))
}

func (propagator b3Propagator) Fields() []string {
	return []string{b3.TraceIDHeader, b3.SpanIDHeader, b3.SampledHeader}
}

func init() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(b3Propagator{
		b3.HTTPFormat{},
	}))
}

func ExtractGinCtx(ctx *gin.Context) context.Context {
	return ExtractHttpRequest(ctx.Request)
}

func ExtractHttpRequest(r *http.Request) context.Context {
	return otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
}

func Inject(ctx context.Context, r *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
}
