package queueProcessor

import (
	"context"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type engineImpl struct {
	config Config
	ProcessingFunction
	gox.CrossFunction
}

func NewEngine(cf gox.CrossFunction, config Config) Engine {
	if util.IsStringEmpty(config.Name) {
		config.Name = uuid.NewString()
	}
	if config.EventBuffer <= 0 {
		config.EventBuffer = 1
	}
	if config.ProcessingFunction == nil {
		config.ProcessingFunction = NewNoOpProcessingFunction()
	}
	return &engineImpl{
		CrossFunction:      cf,
		config:             config,
		ProcessingFunction: config.ProcessingFunction,
	}
}

func (e *engineImpl) StartProcessing(ctx context.Context, rawEventChannel chan RawEvent) (processedEventChannel chan ProcessedEvent) {
	processedEventChannel = make(chan ProcessedEvent, e.config.EventBuffer)
	go func() {
		for {
			select {
			case <-ctx.Done():
				e.Logger().Info("stopping processor", zap.Any("config", e.config), zap.Any("closeEvent", ctx.Err()))
				close(processedEventChannel)
				return
			case rawEvent, ok := <-rawEventChannel:
				if ok {
					processedEvent := e.ProcessEvent(rawEvent)
					if processedEvent.Err == nil {
						processedEventChannel <- processedEvent
					} else {
						e.Logger().Error("failed to process rawEvent", zap.Any("rawEvent", rawEvent))
					}
				} else {
					e.Logger().Info("input raw rawEvent channel is closed: stop processor", zap.Any("name", e.config.Name))
					close(processedEventChannel)
					return
				}
			}
		}
	}()
	return processedEventChannel
}

// Dummy no op processing function
type noOpProcessingFunction struct {
}

func (op *noOpProcessingFunction) ProcessEvent(event RawEvent) ProcessedEvent {
	return ProcessedEvent{
		Data: event.Data,
		Err:  nil,
	}
}
