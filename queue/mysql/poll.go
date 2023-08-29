package queue

import (
	"github.com/devlibx/gox-base/queue"
	"go.uber.org/zap"
	"time"
)

func (t *topRowFinder) Start() {
	go t.internalStart()
}

func (t *topRowFinder) Stop() {
	t.stop = true
}

func (t *topRowFinder) internalStart() {
	for {
		var _time string
		err := t.db.QueryRow("SELECT process_at FROM jobs WHERE tenant=? AND state=? AND job_type=? order by process_at", t.tenant, queue.StatusScheduled, t.jobType).Scan(&_time)
		if err == nil && _time != "" {
			t.smallestProcessAtTime, _ = time.Parse("2006-01-02 15:04:05", _time)
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
			t.logger.Info("[WARN - expected at boot-up] did not find smallest processed at time or if there is not job for given job type", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
		}

		if t.stop {
			t.logger.Info("Stopping top row finder", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
			break
		}
	}
}
