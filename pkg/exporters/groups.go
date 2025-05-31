package exporters

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"netbird-api-exporter/pkg/netbird"
)

// GroupsExporter handles groups-specific metrics collection
type GroupsExporter struct {
	client *netbird.Client

	// Prometheus metrics for groups
	groupsTotal          *prometheus.GaugeVec
	groupPeersCount      *prometheus.GaugeVec
	groupResourcesCount  *prometheus.GaugeVec
	groupInfo            *prometheus.GaugeVec
	groupResourcesByType *prometheus.GaugeVec
	scrapeErrorsTotal    *prometheus.CounterVec
	scrapeDuration       *prometheus.HistogramVec
}

// NewGroupsExporter creates a new groups exporter
func NewGroupsExporter(client *netbird.Client) *GroupsExporter {
	return &GroupsExporter{
		client: client,

		groupsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_groups_total",
				Help: "Total number of NetBird groups",
			},
			[]string{},
		),

		groupPeersCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_group_peers_count",
				Help: "Number of peers in each NetBird group",
			},
			[]string{"group_id", "group_name", "issued"},
		),

		groupResourcesCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_group_resources_count",
				Help: "Number of resources in each NetBird group",
			},
			[]string{"group_id", "group_name", "issued"},
		),

		groupInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_group_info",
				Help: "Information about NetBird groups (always 1)",
			},
			[]string{"group_id", "group_name", "issued"},
		),

		groupResourcesByType: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_group_resources_by_type",
				Help: "Number of resources in each NetBird group by resource type",
			},
			[]string{"group_id", "group_name", "resource_type"},
		),

		scrapeErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "netbird_groups_scrape_errors_total",
				Help: "Total number of errors encountered while scraping groups",
			},
			[]string{"error_type"},
		),

		scrapeDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "netbird_groups_scrape_duration_seconds",
				Help: "Time spent scraping groups from the NetBird API",
			},
			[]string{},
		),
	}
}

// Describe implements prometheus.Collector
func (e *GroupsExporter) Describe(ch chan<- *prometheus.Desc) {
	e.groupsTotal.Describe(ch)
	e.groupPeersCount.Describe(ch)
	e.groupResourcesCount.Describe(ch)
	e.groupInfo.Describe(ch)
	e.groupResourcesByType.Describe(ch)
	e.scrapeErrorsTotal.Describe(ch)
	e.scrapeDuration.Describe(ch)
}

// Collect implements prometheus.Collector
func (e *GroupsExporter) Collect(ch chan<- prometheus.Metric) {
	timer := prometheus.NewTimer(e.scrapeDuration.WithLabelValues())
	defer timer.ObserveDuration()

	// Reset metrics before collecting new values
	e.groupsTotal.Reset()
	e.groupPeersCount.Reset()
	e.groupResourcesCount.Reset()
	e.groupInfo.Reset()
	e.groupResourcesByType.Reset()

	groups, err := e.fetchGroups()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch groups")
		e.scrapeErrorsTotal.WithLabelValues("fetch_groups").Inc()
		return
	}

	e.updateMetrics(groups)

	// Collect all metrics
	e.groupsTotal.Collect(ch)
	e.groupPeersCount.Collect(ch)
	e.groupResourcesCount.Collect(ch)
	e.groupInfo.Collect(ch)
	e.groupResourcesByType.Collect(ch)
	e.scrapeErrorsTotal.Collect(ch)
	e.scrapeDuration.Collect(ch)
}

// fetchGroups retrieves groups from NetBird API
func (e *GroupsExporter) fetchGroups() ([]netbird.Group, error) {
	url := fmt.Sprintf("%s/api/groups", e.client.GetBaseURL())

	logrus.WithField("url", url).Debug("Fetching groups from NetBird API")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", e.client.GetToken()))

	resp, err := e.client.GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logrus.WithError(closeErr).Warn("Failed to close response body for groups")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var groups []netbird.Group
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	logrus.WithField("count", len(groups)).Debug("Successfully fetched groups from API")

	return groups, nil
}

// updateMetrics updates Prometheus metrics based on groups data
func (e *GroupsExporter) updateMetrics(groups []netbird.Group) {
	totalGroups := len(groups)

	// Count resources by type across all groups
	resourceTypeCount := make(map[string]map[string]int) // group_id -> resource_type -> count
	totalPeers := 0
	totalResources := 0
	resourceTypeTotals := make(map[string]int)

	for _, group := range groups {
		groupLabels := []string{group.ID, group.Name, group.Issued}

		// Set basic group metrics
		e.groupPeersCount.WithLabelValues(groupLabels...).Set(float64(group.PeersCount))
		e.groupResourcesCount.WithLabelValues(groupLabels...).Set(float64(group.ResourcesCount))
		e.groupInfo.WithLabelValues(groupLabels...).Set(1)

		// Add to totals
		totalPeers += group.PeersCount
		totalResources += group.ResourcesCount

		// Count resources by type for this group
		if resourceTypeCount[group.ID] == nil {
			resourceTypeCount[group.ID] = make(map[string]int)
		}

		for _, resource := range group.Resources {
			resourceTypeCount[group.ID][resource.Type]++
			resourceTypeTotals[resource.Type]++
		}

		// Set resource type metrics
		for resourceType, count := range resourceTypeCount[group.ID] {
			e.groupResourcesByType.WithLabelValues(group.ID, group.Name, resourceType).Set(float64(count))
		}
	}

	e.groupsTotal.WithLabelValues().Set(float64(totalGroups))

	logrus.WithFields(logrus.Fields{
		"total_groups":          totalGroups,
		"total_peers_in_groups": totalPeers,
		"total_resources":       totalResources,
		"resource_types":        len(resourceTypeTotals),
		"resource_type_counts":  resourceTypeTotals,
	}).Debug("Updated group metrics")
}
