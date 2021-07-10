package util

import "github.com/opentracing/opentracing-go"

func OpentracingLogError(spanName string, err error) {
	if opentracing.GlobalTracer() != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName)
		span.SetTag("error", err)
		span.Finish()
	}
}
