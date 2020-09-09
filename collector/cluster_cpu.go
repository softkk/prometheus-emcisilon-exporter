package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type clusterCPUCollector struct {
	cpuCount  *prometheus.Desc
	cpuIdle   *prometheus.Desc
	cpuUser   *prometheus.Desc
	cpuSys    *prometheus.Desc
	load1min  *prometheus.Desc
	load5min  *prometheus.Desc
	load15min *prometheus.Desc
}

func init() {
	registerCollector("cluster_cpu", defaultEnabled, NewClusterCPUCollector)
}

// NewClusterCPUCollector -
func NewClusterCPUCollector() (Collector, error) {
	return &clusterCPUCollector{
		cpuCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "cpu_count"),
			"Count of number of cpu a node contains.",
			nil, ConstLabels,
		),
		// Cluster average of sytem CPU usage in tenths of a percent
		cpuIdle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "cpu_idle_avg"),
			"Current cpu idle percentage for the node.",
			nil, ConstLabels,
		),
		cpuUser: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "cpu_user_avg"),
			"Current cpu busy percentage for user mode represented in 0.0-1.0.",
			nil, ConstLabels,
		),
		cpuSys: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "cpu_sys_avg"),
			"Current cpu busy percentage for sys mode represented in 0.0-1.0.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *clusterCPUCollector) Update(ch chan<- prometheus.Metric) error {
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.cpuCount] = "cluster.cpu.count"
	keyMap[c.cpuIdle] = "cluster.cpu.idle.avg"
	keyMap[c.cpuUser] = "cluster.cpu.user.avg"
	keyMap[c.cpuSys] = "cluster.cpu.sys.avg"

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
				var val float64
				if strings.Contains(statKey, "cpu") {
					if !(strings.Contains(statKey, "count")) {
						val = stat.Value / 10 // 原本是1/10, 除10變成百分之一
					} else {
						val = stat.Value
					}
				}

				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, val)
			}
		}
	}

	if errCount != 0 {
		err := fmt.Errorf("There where %v errors", errCount)
		return err
	}
	return nil
}
