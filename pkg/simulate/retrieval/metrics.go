package retrieval

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
)

type metrics struct {
	uploadedCounter       *prometheus.CounterVec
	notUploadedCounter    *prometheus.CounterVec
	uploadTimeGauge       *prometheus.GaugeVec
	uploadTimeHistogram   prometheus.Histogram
	downloadedCounter     *prometheus.CounterVec
	notDownloadedCounter  *prometheus.CounterVec
	downloadTimeGauge     *prometheus.GaugeVec
	downloadTimeHistogram prometheus.Histogram
	retrievedCounter      *prometheus.CounterVec
	notRetrievedCounter   *prometheus.CounterVec
	syncedCounter         *prometheus.CounterVec
	notSyncedCounter      *prometheus.CounterVec
	syncTagsTimeGauge     *prometheus.GaugeVec
	syncTagsTimeHistogram prometheus.Histogram
}

func newMetrics(runID string, pusher *push.Pusher) metrics {
	namespace := "beekeeper"
	subsystem := "simulation_retrieval"

	addCollector := func(c prometheus.Collector) {
		if pusher != nil {
			pusher.Collector(c)
		}
	}

	uploadedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_uploaded_count",
			Help: "Number of uploaded chunks.",
		},
		[]string{"node"},
	)
	addCollector(uploadedCounter)

	notUploadedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_not_uploaded_count",
			Help: "Number of not uploaded chunks.",
		},
		[]string{"node"},
	)
	addCollector(notUploadedCounter)

	uploadTimeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunk_upload_duration_seconds",
			Help: "Chunk upload duration Gauge.",
		},
		[]string{"node", "chunk"},
	)
	addCollector(uploadTimeGauge)

	uploadTimeHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name:    "chunk_upload_seconds",
			Help:    "Chunk upload duration Histogram.",
			Buckets: []float64{0.01, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)
	addCollector(uploadTimeHistogram)

	downloadedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_downloaded_count",
			Help: "Number of downloaded chunks.",
		},
		[]string{"node"},
	)
	addCollector(downloadedCounter)

	notDownloadedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_not_downloaded_count",
			Help: "Number of chunks that has not been downloaded.",
		},
		[]string{"node"},
	)
	addCollector(notDownloadedCounter)

	downloadTimeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunk_download_duration_seconds",
			Help: "Chunk download duration Gauge.",
		},
		[]string{"node", "chunk"},
	)
	addCollector(downloadTimeGauge)

	downloadTimeHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name:    "chunk_download_seconds",
			Help:    "Chunk download duration Histogram.",
			Buckets: []float64{0.01, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)
	addCollector(downloadTimeHistogram)

	retrievedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_retrieved_count",
			Help: "Number of chunks that has been retrieved.",
		},
		[]string{"node"},
	)
	addCollector(retrievedCounter)

	notRetrievedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "chunks_not_retrieved_count",
			Help: "Number of chunks that has not been retrieved.",
		},
		[]string{"node"},
	)
	addCollector(notRetrievedCounter)

	syncedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "tags_synced_count",
			Help: "Number of synced tags.",
		},
		[]string{"node"},
	)
	addCollector(syncedCounter)

	notSyncedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "tags_not_synced_count",
			Help: "Number of not synced tags.",
		},
		[]string{"node"},
	)
	addCollector(notSyncedCounter)

	syncTagsTimeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name: "tags_sync_duration_seconds",
			Help: "Tags sync duration Gauge.",
		},
		[]string{"node", "chunk"},
	)
	addCollector(syncTagsTimeGauge)

	syncTagsTimeHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			ConstLabels: prometheus.Labels{
				"run": runID,
			},
			Name:    "tags_sync_seconds",
			Help:    "Tags sync duration Histogram.",
			Buckets: []float64{0.01, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)
	addCollector(syncTagsTimeHistogram)

	if pusher != nil {
		pusher.Format(expfmt.FmtText)
	}

	return metrics{
		uploadedCounter:       uploadedCounter,
		notUploadedCounter:    notUploadedCounter,
		uploadTimeGauge:       uploadTimeGauge,
		uploadTimeHistogram:   uploadTimeHistogram,
		downloadedCounter:     downloadedCounter,
		notDownloadedCounter:  notDownloadedCounter,
		downloadTimeGauge:     downloadTimeGauge,
		downloadTimeHistogram: downloadTimeHistogram,
		retrievedCounter:      retrievedCounter,
		notRetrievedCounter:   notRetrievedCounter,
		syncedCounter:         syncedCounter,
		notSyncedCounter:      notSyncedCounter,
		syncTagsTimeGauge:     syncTagsTimeGauge,
		syncTagsTimeHistogram: syncTagsTimeHistogram,
	}
}
