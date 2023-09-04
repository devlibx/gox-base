package queue

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"sync"
	"time"
)

type IdGenerator interface {
	GenerateId(input interface{}) string
}

type RandomUuidIdGenerator struct {
}

func (r RandomUuidIdGenerator) GenerateId(input interface{}) string {
	return uuid.NewString()
}

func NewRandomUuidIdGenerator() (IdGenerator, error) {
	return &RandomUuidIdGenerator{}, nil
}

type TimeBasedIdGenerator struct {
	entropy *rand.Rand
	m       *sync.Mutex
}

func (t *TimeBasedIdGenerator) GenerateId(input interface{}) string {
	if inTime, ok := input.(time.Time); ok {
		t.m.Lock()
		defer t.m.Unlock()
		ms := ulid.Timestamp(inTime)
		if r, err := ulid.New(ms, t.entropy); err == nil {
			return r.String()
		} else {
			fmt.Println("[WARN] failed to generate ulid from TimeBasedIdGenerator", err)
			return uuid.NewString()
		}
	} else {
		fmt.Println("[WARN] failed to generate ulid from TimeBasedIdGenerator because input is not time.Time. input=", input)
		return uuid.NewString()
	}
}

func NewTimeBasedIdGenerator() (IdGenerator, error) {
	t := &TimeBasedIdGenerator{
		entropy: rand.New(rand.NewSource(time.Now().UnixNano())),
		m:       &sync.Mutex{},
	}
	return t, nil
}

// RetryBackoffAlgo will help to schedule next retry
type RetryBackoffAlgo interface {
	NextRetryAfter(attempt int, maxExecution int) (time.Duration, error)
}
