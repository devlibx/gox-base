package queueProcessor

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/v2/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkEngine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctrl := gomock.NewController(b)
		cf := test.BuildMockCfB(b, ctrl)

		engine := NewEngine(cf, Config{
			Script:             "",
			EventBuffer:        1,
			ProcessingFunction: NewNoOpProcessingFunction(),
		})

		ctx, _ := context.WithCancel(context.TODO())
		rawEventChannel := make(chan RawEvent, 1000)
		processedEventChannel := engine.StartProcessing(ctx, rawEventChannel)

		go func() {
			rawEventChannel <- RawEvent{Data: map[string]interface{}{"in": 1}}
			rawEventChannel <- RawEvent{Data: map[string]interface{}{"in": 2}}
			close(rawEventChannel)
		}()

		count := 0
		for _ = range processedEventChannel {
			count++
		}
		assert.Equal(b, 2, count)
	}
}

func TestEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	cf := test.BuildMockCf(t, ctrl)

	engine := NewEngine(cf, Config{
		Script:             "",
		EventBuffer:        1,
		ProcessingFunction: NewNoOpProcessingFunction(),
	})

	ctx, _ := context.WithCancel(context.TODO())
	rawEventChannel := make(chan RawEvent, 100)
	processedEventChannel := engine.StartProcessing(ctx, rawEventChannel)

	go func() {
		rawEventChannel <- RawEvent{Data: map[string]interface{}{"in": 1}}
		rawEventChannel <- RawEvent{Data: map[string]interface{}{"in": 2}}
		close(rawEventChannel)
	}()

	count := 0
	for event := range processedEventChannel {
		count++
		fmt.Println(event)
	}
	assert.Equal(t, 2, count)
}
