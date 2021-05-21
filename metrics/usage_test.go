package metrics

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestTypes_UsageSample(t *testing.T) {
	if true {
		return
	}
	var r Reporter
	var scope Scope

	counter := scope.Tagged(map[string]string{
		"foo": "bar",
	}).Counter("test_counter")

	ti := scope.Tagged(map[string]string{"name": "harish"}).Timer("harish_timer")

	gauge := scope.Tagged(map[string]string{
		"foo": "baz",
	}).Gauge("test_gauge")

	timer := scope.Tagged(map[string]string{
		"foo": "qux",
	}).Timer("test_timer_summary")

	histogram := scope.Tagged(map[string]string{
		"foo": "quk",
	}).Histogram("test_histogram", DefaultBuckets)

	go func() {
		for {
			tr := ti.Start()
			counter.Inc(1)
			time.Sleep(time.Second)
			tr.Stop()
		}
	}()

	go func() {
		for {
			gauge.Update(rand.Float64() * 1000)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tsw := timer.Start()
			hsw := histogram.Start()
			time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
			tsw.Stop()
			hsw.Stop()
		}
	}()

	http.Handle("/metrics", r.HTTPHandler())
	fmt.Printf("Serving :8080/metrics\n")
	fmt.Printf("%v\n", http.ListenAndServe(":8089", nil))
	select {}
}
