package exporters

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"netbird-api-exporter/pkg/netbird"
)

// NetBirdExporter represents the main Prometheus exporter for NetBird APIs
type NetBirdExporter struct {
	client           *netbird.Client
	peersExporter    *PeersExporter
	groupsExporter   *GroupsExporter
	usersExporter    *UsersExporter
	dnsExporter      *DNSExporter
	networksExporter *NetworksExporter

	// Common metrics
	scrapeDuration prometheus.Histogram
	scrapeErrors   prometheus.Counter
}

// NewNetBirdExporter creates a new NetBird exporter with all sub-exporters
func NewNetBirdExporter(baseURL, token string) *NetBirdExporter {
	client := netbird.NewClient(baseURL, token)

	return &NetBirdExporter{
		client:           client,
		peersExporter:    NewPeersExporter(client),
		groupsExporter:   NewGroupsExporter(client),
		usersExporter:    NewUsersExporter(client),
		dnsExporter:      NewDNSExporter(client),
		networksExporter: NewNetworksExporter(client),

		scrapeDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "netbird_exporter_scrape_duration_seconds",
				Help: "Time spent scraping NetBird API",
			},
		),

		scrapeErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "netbird_exporter_scrape_errors_total",
				Help: "Total number of scrape errors",
			},
		),
	}
}

// Describe implements prometheus.Collector
func (e *NetBirdExporter) Describe(ch chan<- *prometheus.Desc) {
	e.peersExporter.Describe(ch)
	e.groupsExporter.Describe(ch)
	e.usersExporter.Describe(ch)
	e.dnsExporter.Describe(ch)
	e.networksExporter.Describe(ch)
	e.scrapeDuration.Describe(ch)
	e.scrapeErrors.Describe(ch)
}

// Collect implements prometheus.Collector
func (e *NetBirdExporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	defer func() {
		e.scrapeDuration.Observe(time.Since(start).Seconds())
		e.scrapeDuration.Collect(ch)
		e.scrapeErrors.Collect(ch)
	}()

	// Collect from all sub-exporters
	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during peers collection")
				e.scrapeErrors.Inc()
			}
		}()
		e.peersExporter.Collect(ch)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during groups collection")
				e.scrapeErrors.Inc()
			}
		}()
		e.groupsExporter.Collect(ch)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during users collection")
				e.scrapeErrors.Inc()
			}
		}()
		e.usersExporter.Collect(ch)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during dns collection")
				e.scrapeErrors.Inc()
			}
		}()
		e.dnsExporter.Collect(ch)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during networks collection")
				e.scrapeErrors.Inc()
			}
		}()
		e.networksExporter.Collect(ch)
	}()

	// Future exporters can be added here like:
	// e.policiesExporter.Collect(ch)
}
