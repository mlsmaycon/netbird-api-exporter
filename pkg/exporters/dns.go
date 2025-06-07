package exporters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"netbird-api-exporter/pkg/netbird"
)

// DNSExporter handles DNS-specific metrics collection
type DNSExporter struct {
	client *netbird.Client

	// Prometheus metrics
	nameserverGroupsTotal   *prometheus.GaugeVec
	nameserverGroupsEnabled *prometheus.GaugeVec
	nameserverGroupsPrimary *prometheus.GaugeVec
	nameserverGroupDomains  *prometheus.GaugeVec
	nameserversTotal        *prometheus.GaugeVec
	nameserversByType       *prometheus.GaugeVec
	nameserversByPort       *prometheus.GaugeVec
	dnsManagementDisabled   *prometheus.GaugeVec
}

// NewDNSExporter creates a new DNS exporter
func NewDNSExporter(client *netbird.Client) *DNSExporter {
	return &DNSExporter{
		client: client,

		nameserverGroupsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameserver_groups",
				Help: "Total number of NetBird nameserver groups",
			},
			[]string{},
		),

		nameserverGroupsEnabled: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameserver_groups_enabled",
				Help: "Number of enabled NetBird nameserver groups",
			},
			[]string{"enabled"},
		),

		nameserverGroupsPrimary: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameserver_groups_primary",
				Help: "Number of primary NetBird nameserver groups",
			},
			[]string{"primary"},
		),

		nameserverGroupDomains: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameserver_group_domains_count",
				Help: "Number of domains configured in each nameserver group",
			},
			[]string{"group_id", "group_name"},
		),

		nameserversTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameservers",
				Help: "Total number of nameservers across all groups",
			},
			[]string{"group_id", "group_name"},
		),

		nameserversByType: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameservers_by_type",
				Help: "Number of nameservers by type (UDP/TCP)",
			},
			[]string{"ns_type"},
		),

		nameserversByPort: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_nameservers_by_port",
				Help: "Number of nameservers by port",
			},
			[]string{"port"},
		),

		dnsManagementDisabled: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "netbird_dns_management_disabled_groups_count",
				Help: "Number of groups with DNS management disabled",
			},
			[]string{},
		),
	}
}

// Describe implements prometheus.Collector
func (e *DNSExporter) Describe(ch chan<- *prometheus.Desc) {
	e.nameserverGroupsTotal.Describe(ch)
	e.nameserverGroupsEnabled.Describe(ch)
	e.nameserverGroupsPrimary.Describe(ch)
	e.nameserverGroupDomains.Describe(ch)
	e.nameserversTotal.Describe(ch)
	e.nameserversByType.Describe(ch)
	e.nameserversByPort.Describe(ch)
	e.dnsManagementDisabled.Describe(ch)
}

// Collect implements prometheus.Collector
func (e *DNSExporter) Collect(ch chan<- prometheus.Metric) {
	// Reset metrics before collecting new values
	e.nameserverGroupsTotal.Reset()
	e.nameserverGroupsEnabled.Reset()
	e.nameserverGroupsPrimary.Reset()
	e.nameserverGroupDomains.Reset()
	e.nameserversTotal.Reset()
	e.nameserversByType.Reset()
	e.nameserversByPort.Reset()
	e.dnsManagementDisabled.Reset()

	// Fetch nameserver groups
	nameserverGroups, err := e.fetchNameserverGroups()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch nameserver groups")
	} else {
		e.updateNameserverMetrics(nameserverGroups)
	}

	// Fetch DNS settings
	dnsSettings, err := e.fetchDNSSettings()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch DNS settings")
	} else {
		e.updateDNSSettingsMetrics(dnsSettings)
	}

	// Collect all metrics
	e.nameserverGroupsTotal.Collect(ch)
	e.nameserverGroupsEnabled.Collect(ch)
	e.nameserverGroupsPrimary.Collect(ch)
	e.nameserverGroupDomains.Collect(ch)
	e.nameserversTotal.Collect(ch)
	e.nameserversByType.Collect(ch)
	e.nameserversByPort.Collect(ch)
	e.dnsManagementDisabled.Collect(ch)
}

// fetchNameserverGroups retrieves nameserver groups from NetBird API
func (e *DNSExporter) fetchNameserverGroups() ([]netbird.NameserverGroup, error) {
	url := fmt.Sprintf("%s/api/dns/nameservers", e.client.GetBaseURL())

	logrus.WithField("url", url).Debug("Fetching nameserver groups from NetBird API")

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
			logrus.WithError(closeErr).Warn("Failed to close response body for nameserver groups")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var nameserverGroups []netbird.NameserverGroup
	if err := json.NewDecoder(resp.Body).Decode(&nameserverGroups); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	logrus.WithField("count", len(nameserverGroups)).Debug("Successfully fetched nameserver groups from API")

	return nameserverGroups, nil
}

// fetchDNSSettings retrieves DNS settings from NetBird API
func (e *DNSExporter) fetchDNSSettings() (*netbird.DNSSettings, error) {
	url := fmt.Sprintf("%s/api/dns/settings", e.client.GetBaseURL())

	logrus.WithField("url", url).Debug("Fetching DNS settings from NetBird API")

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
			logrus.WithError(closeErr).Warn("Failed to close response body for DNS settings")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var dnsSettings netbird.DNSSettings
	if err := json.NewDecoder(resp.Body).Decode(&dnsSettings); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	logrus.Debug("Successfully fetched DNS settings from API")

	return &dnsSettings, nil
}

// updateNameserverMetrics updates Prometheus metrics based on nameserver group data
func (e *DNSExporter) updateNameserverMetrics(nameserverGroups []netbird.NameserverGroup) {
	// Count totals
	totalGroups := len(nameserverGroups)
	e.nameserverGroupsTotal.WithLabelValues().Set(float64(totalGroups))

	// Count by status
	enabledCounts := make(map[bool]int)
	primaryCounts := make(map[bool]int)
	typeCounter := make(map[string]int)
	portCounter := make(map[string]int)

	for _, group := range nameserverGroups {
		// Count enabled/disabled
		enabledCounts[group.Enabled]++

		// Count primary/secondary
		primaryCounts[group.Primary]++

		// Count domains per group
		e.nameserverGroupDomains.WithLabelValues(group.ID, group.Name).Set(float64(len(group.Domains)))

		// Count nameservers per group
		e.nameserversTotal.WithLabelValues(group.ID, group.Name).Set(float64(len(group.Nameservers)))

		// Count nameserver types and ports
		for _, ns := range group.Nameservers {
			typeCounter[ns.NSType]++
			portCounter[strconv.Itoa(ns.Port)]++
		}
	}

	// Set enabled/disabled metrics
	for enabled, count := range enabledCounts {
		e.nameserverGroupsEnabled.WithLabelValues(strconv.FormatBool(enabled)).Set(float64(count))
	}

	// Set primary/secondary metrics
	for primary, count := range primaryCounts {
		e.nameserverGroupsPrimary.WithLabelValues(strconv.FormatBool(primary)).Set(float64(count))
	}

	// Set nameserver type metrics
	for nsType, count := range typeCounter {
		e.nameserversByType.WithLabelValues(nsType).Set(float64(count))
	}

	// Set nameserver port metrics
	for port, count := range portCounter {
		e.nameserversByPort.WithLabelValues(port).Set(float64(count))
	}

	logrus.WithFields(logrus.Fields{
		"total_groups":    totalGroups,
		"enabled_groups":  enabledCounts[true],
		"disabled_groups": enabledCounts[false],
		"primary_groups":  primaryCounts[true],
	}).Debug("Updated nameserver metrics")
}

// updateDNSSettingsMetrics updates Prometheus metrics based on DNS settings data
func (e *DNSExporter) updateDNSSettingsMetrics(dnsSettings *netbird.DNSSettings) {
	// Count disabled management groups
	disabledCount := len(dnsSettings.Items.DisabledManagementGroups)
	e.dnsManagementDisabled.WithLabelValues().Set(float64(disabledCount))

	logrus.WithField("disabled_management_groups", disabledCount).Debug("Updated DNS settings metrics")
}
