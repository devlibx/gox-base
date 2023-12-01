package metrics

import "context"

// Publishable is an interface which can be used to publish using publisher
//
// Note - you can publish metric using metric.Scope -> it is mostly meant for publishing Prometheus, StatsD etc
// However there are times when we want to send high dim, high cardinality data to kafka or other system (as a metric)
// These event can then pe pushed to Druid, Data lake and gives business event visibility
type Publishable interface {

	// Payload will return key and value which will be used to publish the object
	// key - this is a partition id e.g. if you use kafka has a backend to publish, this will be partition key
	// value - this is the payload which is sent as a data
	Payload() (key interface{}, value interface{})
}

// Publisher is an interface which can be used to publish a object
type Publisher interface {

	// Publish will publish the object
	Publish(ctx context.Context, p Publishable) error

	// SilentPublish will publish the object but will not return any error
	// Since these objects can be lost or ignored, you may choose to ignore errors
	SilentPublish(ctx context.Context, p Publishable)
}

// noOpPublisher is a publisher which does nothing
type noOpPublisher struct {
}

func (n *noOpPublisher) Publish(ctx context.Context, p Publishable) error {
	return nil
}

func (n *noOpPublisher) SilentPublish(ctx context.Context, p Publishable) {
}

func NewNoOpPublisher() Publisher {
	return &noOpPublisher{}
}
