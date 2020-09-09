package collector

import (
	"fmt"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type nodeCapacityCollector struct {
	nodeBytesTotal   *prometheus.Desc
	nodeBytesUsed    *prometheus.Desc
	nodeBytesFree    *prometheus.Desc
	nodeBytesInRate  *prometheus.Desc
	nodeBytesOutRate *prometheus.Desc
}

func init() {
	registerCollector("node_capacity", defaultEnabled, NewNodeCapacityCollector)
}

//NewNodeCapacityCollector -
func NewNodeCapacityCollector() (Collector, error) {
	return &nodeCapacityCollector{
		nodeBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "bytes_total"),
			"Current ifs filesystem capacity on the node in bytes.",
			[]string{"node"}, ConstLabels,
		),
		nodeBytesUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "bytes_used"),
			"Current ifs filesystem capacity used on the node in bytes.",
			[]string{"node"}, ConstLabels,
		),
		nodeBytesFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "bytes_free"),
			"Current ifs filesystem capacity free on the node in bytes.",
			[]string{"node"}, ConstLabels,
		),
		nodeBytesInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "bytes_in_rate"),
			"Input rate to /ifs from the node (bytes/sec)",
			[]string{"node"}, ConstLabels,
		),
		nodeBytesOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "bytes_out_rate"),
			"Output rate to /ifs from the node (bytes/sec).",
			[]string{"node"}, ConstLabels,
		),
	}, nil
}

//Update -
func (c *nodeCapacityCollector) Update(ch chan<- prometheus.Metric) error {
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.nodeBytesTotal] = "node.ifs.bytes.total"
	keyMap[c.nodeBytesUsed] = "node.ifs.bytes.used"
	keyMap[c.nodeBytesFree] = "node.ifs.bytes.free"
	keyMap[c.nodeBytesInRate] = "node.ifs.bytes.in.rate"
	keyMap[c.nodeBytesOutRate] = "node.ifs.bytes.out.rate"
	//  node.ifs.bytes.in.rate.max
	//  node.ifs.bytes.out.rate.max

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
				node := fmt.Sprintf("%v", stat.Devid)
				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, stat.Value, node)
			}
		}

	}

	if errCount != 0 {
		err := fmt.Errorf("There where %d errors", errCount)
		return err
	}
	return nil
}
