package collector

import (
	"fmt"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type clusterNetworkCollector struct {
	netBytesInRate   *prometheus.Desc
	netBytesOutRate  *prometheus.Desc
	netErrorsInRate  *prometheus.Desc
	netErrorsOutRate *prometheus.Desc
}

func init() {
	registerCollector("cluster_network", defaultEnabled, NewClusterNetworkCollector)
}

// NewClusterNetworkCollector -
func NewClusterNetworkCollector() (Collector, error) {
	return &clusterNetworkCollector{
		netBytesInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "net_ext_bytes_in_rate"),
			"Current network bytes per second in rate from all external interfaces.",
			nil, ConstLabels,
		),
		netBytesOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "net_ext_bytes_out_rate"),
			"Current network bytes out rate from external interfaces.",
			nil, ConstLabels,
		),
		netErrorsInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "net_ext_errors_in_rate"),
			"Input errors per second for a node's external interfaces.",
			nil, ConstLabels,
		),
		netErrorsOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "net_ext_errors_out_rate"),
			"Output errors per seccond for a node's external interfaces.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *clusterNetworkCollector) Update(ch chan<- prometheus.Metric) error {
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.netBytesInRate] = "cluster.net.ext.bytes.in.rate"
	keyMap[c.netBytesOutRate] = "cluster.net.ext.bytes.out.rate"
	keyMap[c.netErrorsInRate] = "cluster.net.ext.errors.in.rate"
	keyMap[c.netErrorsOutRate] = "cluster.net.ext.errors.out.rate"

	for promStat, statKey := range keyMap {
		begin := time.Now()
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		duration := time.Since(begin)
		ch <- prometheus.MustNewConstMetric(statsEngineCallDuration, prometheus.GaugeValue, duration.Seconds(), statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			ch <- prometheus.MustNewConstMetric(statsEngineCallFailure, prometheus.GaugeValue, 1, statKey)
			errCount++
		} else {
			ch <- prometheus.MustNewConstMetric(statsEngineCallFailure, prometheus.GaugeValue, 0, statKey)
			for _, stat := range resp.Stats {
				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, stat.Value)
			}
		}
	}
	if errCount != 0 {
		err := fmt.Errorf("There where %d errors", errCount)
		return err
	}
	return nil
}
