package gox

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// EnableInterceptorInDefaultTimeService is a global flag to enable/disable the interceptor logic in DefaultTimeService.
// When enabled, DefaultTimeService.CurrentTime will check for any registered interceptors before returning the current time.
// This is typically used for testing or when time needs to be controlled/mocked at a global level.
var EnableInterceptorInDefaultTimeService = false

// InterceptorDefaultTimeServiceNowFunction is a function type that can be used to intercept and provide a custom time.
// It takes a context and an optional input, and returns a boolean (indicating if the interceptor handled the request) and the time.
type InterceptorDefaultTimeServiceNowFunction func(ctx context.Context, input any) (bool, time.Time)

// Global map and lock for registered interceptors.
var mapInterceptorDefaultTimeServiceNowFunction map[string]InterceptorDefaultTimeServiceNowFunction
var mapInterceptorDefaultTimeServiceNowFunctionLock *sync.RWMutex

func init() {
	mapInterceptorDefaultTimeServiceNowFunction = map[string]InterceptorDefaultTimeServiceNowFunction{}
	mapInterceptorDefaultTimeServiceNowFunctionLock = &sync.RWMutex{}
}

// InterceptableTimeService defines an interface for a time service that supports interception.
type InterceptableTimeService interface {
	// CurrentTime returns the current time, potentially intercepted by registered functions if EnableInterceptorInDefaultTimeService is true.
	CurrentTime(ctx context.Context, input any) time.Time
}

// DefaultTimeService provides a standard implementation of the TimeService interface.
type DefaultTimeService struct {
	TimeService
}

// Now returns the current wall-clock time using time.Now().
func (t *DefaultTimeService) Now() time.Time {
	return time.Now()
}

// CurrentTime returns the current time. If EnableInterceptorInDefaultTimeService is true, it iterates through
// all registered interceptors and returns the time from the first one that signals it has handled the request (ok=true).
// If no interceptor handles the request or if the feature is disabled, it returns the current wall-clock time.
func (t *DefaultTimeService) CurrentTime(ctx context.Context, input any) time.Time {
	if EnableInterceptorInDefaultTimeService {
		mapInterceptorDefaultTimeServiceNowFunctionLock.RLock()
		for k, v := range mapInterceptorDefaultTimeServiceNowFunction {
			if ok, t := v(ctx, input); ok {
				slog.Warn("*** returning time from interceptor ***", "id", k)
				return t
			}
		}
		defer mapInterceptorDefaultTimeServiceNowFunctionLock.RUnlock()
	}
	return time.Now()
}

// Sleep pauses the current goroutine for the given duration using time.Sleep().
func (t *DefaultTimeService) Sleep(d time.Duration) {
	time.Sleep(d)
}

// NormalizeDuration returns the duration as-is for the default implementation.
func (t *DefaultTimeService) NormalizeDuration(d time.Duration) time.Duration {
	return d
}

// RegisterInterceptorDefaultTimeServiceNowFunction adds a new interceptor with the given ID.
// The interceptor will only be added and used if EnableInterceptorInDefaultTimeService is true.
func RegisterInterceptorDefaultTimeServiceNowFunction(id string, f InterceptorDefaultTimeServiceNowFunction) {
	if EnableInterceptorInDefaultTimeService {
		mapInterceptorDefaultTimeServiceNowFunctionLock.Lock()
		mapInterceptorDefaultTimeServiceNowFunction[id] = f
		mapInterceptorDefaultTimeServiceNowFunctionLock.Unlock()
	}
}

// UnregisterInterceptorDefaultTimeServiceNowFunction removes the interceptor with the given ID.
// It only performs the removal if EnableInterceptorInDefaultTimeService is true.
func UnregisterInterceptorDefaultTimeServiceNowFunction(id string) {
	if EnableInterceptorInDefaultTimeService {
		mapInterceptorDefaultTimeServiceNowFunctionLock.Lock()
		delete(mapInterceptorDefaultTimeServiceNowFunction, id)
		mapInterceptorDefaultTimeServiceNowFunctionLock.Unlock()
	}
}
