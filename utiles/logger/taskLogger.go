package logger

import (
	"sync"
	"time"
)

// ----------------|

type TaskLogger struct {
	points map[string]point
	mu     sync.Mutex
}

type point struct {
	duration time.Duration
	count    int
}

// ----------------|

type TaskStat struct {
	Count     int
	TotalTime time.Duration
	AvgTime   time.Duration
}

type TaskStats map[string]TaskStat

// ----------------|

// PushPoint - log current task
// NOTE:
// 	- call the method with defer:
//		defer lg.PushPoint("name", time.Now())
func (lg *TaskLogger) PushPoint(name string, start time.Time) {
	duration := time.Since(start)

	lg.mu.Lock()
	defer lg.mu.Unlock()
	if p, exists := lg.points[name]; !exists {
		lg.points[name] = point{
			duration: duration,
			count:    1,
		}
	} else {
		p.duration += duration
		p.count++
	}
}

// CreateStatistics - generates statistics about the execution time of tasks
func (lg *TaskLogger) CreateStatistics() TaskStats {
	stats := TaskStats{}

	lg.mu.Lock()
	defer lg.mu.Unlock()
	for name, p := range lg.points {
		stats[name] = TaskStat{
			Count:     p.count,
			TotalTime: p.duration,
			AvgTime:   time.Duration(int64(p.duration) / int64(p.count)),
		}
	}
	return stats
}
