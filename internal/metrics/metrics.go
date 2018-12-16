package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/node_exporter"
)

type Profiler struct {
	HitsStats   prometheus.CounterVec
	ActiveRooms prometheus.Gauge
}

func Construct() *Profiler {
	p := Profiler{
		HitsStats: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hits_statuses_total",
			Help: "Number of hits and statuses.",
		},
			[]string{"code", "status"},
		),

		ActiveRooms: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_rooms",
			Help: "Number of active rooms.",
		}),
	}

	prometheus.MustRegister(p.HitsStats)
	prometheus.MustRegister(p.ActiveRooms)

	return &p
}
