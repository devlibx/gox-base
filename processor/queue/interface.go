package processor

import (
	"context"
	"github.com/devlibx/gox-base"
)

// Converts raw event to processed event
type ProcessingFunction interface {
	ProcessEvent(event RawEvent) ProcessedEvent
}

type Config struct {
	Name               string
	Script             string
	EventBuffer        int
	ProcessingFunction ProcessingFunction
}

type RawEvent struct {
	Data gox.StringObjectMap
}

type ProcessedEvent struct {
	Data gox.StringObjectMap
	Err  error
}

// Engine to process the input event
type Engine interface {

	// Start processing of the events coming in rawEventChannel
	// Output events will be in out channel
	// If you want to stop he work then cancel the input context. NOTE - this will stop the processing immediately
	// which may cause event loss which yoy sent in the rawEventChannel
	// Of you just close the rawEventChannel then the engine will stop once all the events are processed
	StartProcessing(ctx context.Context, rawEventChannel chan RawEvent) (processedEventChannel chan ProcessedEvent)
}

// A processing function which does not do anything
func NewNoOpProcessingFunction() ProcessingFunction {
	return &noOpProcessingFunction{}
}
