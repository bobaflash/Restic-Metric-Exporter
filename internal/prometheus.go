package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	SnapshotTimestamp      *prometheus.GaugeVec
	SnapshotCount          *prometheus.GaugeVec
	SnapshotTotalSize      *prometheus.GaugeVec
	SnapshotTotalFileCount *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		SnapshotTimestamp: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "restic_snapshot_timestamp",
			Help: "The last time a snapshot was written",
		}, []string{"repo_name", "hostname"}),

		SnapshotCount: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "restic_snapshot_count",
			Help: "The amount of snapshots in a repository",
		}, []string{"repo_name"}),

		SnapshotTotalSize: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "restic_snapshot_total_size",
			Help: "The size of all snapshots in bytes within a repository",
		}, []string{"repo_name"}),

		SnapshotTotalFileCount: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "restic_snapshot_total_files",
			Help: "The amount of all files in within a repository",
		}, []string{"repo_name"}),
	}
	return m
}
