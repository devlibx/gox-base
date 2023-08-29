package queue

import (
	"github.com/devlibx/gox-base/queue"
	"go.uber.org/zap"
	"time"
)

func (s *queueImpl) topRowScanner(jobType string) {
	for {
		var t string
		err := s.db.QueryRow("SELECT process_at FROM jobs WHERE state=? AND job_type=? order by process_at asc", queue.StatusScheduled, jobType).Scan(&t)
		if err == nil {
			// s.smallestProcessedAt, err = time.Parse("2006-01-02 15:04:05", t)
			s.smallestProcessedAt[jobType], err = time.Parse("2006-01-02 15:04:05", t)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				s.logger.Error("[WARN] failed to find smallest processed at time", zap.String("t=", t), zap.Error(err))
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			time.Sleep(1 * time.Second)
			s.logger.Warn("[WARN - expected at boot-up] did not find smallest processed at time", zap.String("t=", t), zap.Error(err))
		}
	}
}
