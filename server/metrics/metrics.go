package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/node_exporter"
)

type Profiler struct {
	Hits        prometheus.Counter
	ActiveRooms prometheus.Gauge
	Routines    prometheus.Gauge
	CPUUsage    prometheus.Gauge
	MemUsage    prometheus.Gauge
	DiskUsage   prometheus.Gauge
}

func Construct() *Profiler {
	p := Profiler{
		Hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "hits_statuses_total",
			Help: "Number of hits and statuses.",
		}),

		ActiveRooms: prometheus.NewGauge(prometheus.Gau
			Help: "Number of active rooms.",
		}),

		CPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "Info on CPU usage.",
		}),

		MemUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "mem_usage",
			Help: "Info on mem usage.",
		}),

		DiskUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "disk_usage",
			Help: "Info on disk usage.",
		}),
	}

	prometheus.MustRegister(p.Hits)
	prometheus.MustRegister(p.ActiveRooms)
	prometheus.MustRegister(p.CPUUsage)
	prometheus.MustRegister(p.MemUsage)
	prometheus.MustRegister(p.DiskUsage)

	return &p
}
