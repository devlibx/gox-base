package util

import "github.com/opentracing/opentracing-go"

func OpentracingLogError(spanName string, err error) {
	if opentracing.GlobalTracer() != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName)
		span.SetTag("error", err)
		span.Finish()
	}
}

func OpentracingLogError1(spanName string, err error, key string, value interface{}) {
	if opentracing.GlobalTracer() != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName)
		span.SetTag("error", err)
		span.SetTag(key, value)
		span.Finish()
	}
}

func OpentracingLogError2(spanName string, err error, key string, value interface{}, key1 string, value1 interface{}) {
	if opentracing.GlobalTracer() != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName)
		span.SetTag("error", err)
		span.SetTag(key, value)
		span.SetTag(key1, value1)
		span.Finish()
	}
}

func OpentracingLogError3(spanName string, err error, key string, value interface{}, key1 string, value1 interface{}, key2 string, value2 interface{}) {
	if opentracing.GlobalTracer() != nil {
		span := opentracing.GlobalTracer().StartSpan(spanName)
		span.SetTag("error", err)
		span.SetTag(key, value)
		span.SetTag(key1, value1)
		span.SetTag(key2, value2)
		span.Finish()
	}
}