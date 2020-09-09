package collector

import (
	"strconv"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/cache"
	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type fsaCollector struct {
	clusterFSAAdsCount           *prometheus.Desc
	clusterFSADirCount           *prometheus.Desc
	clusterFSAFileCount          *prometheus.Desc
	clusterFSAHasSubdirs         *prometheus.Desc
	clusterFSALogSizeSum         *prometheus.Desc
	clusterFSALogSizeSumOverflow *prometheus.Desc
	clusterFSAOtherCount         *prometheus.Desc
	clusterFSAPhysSizeSum        *prometheus.Desc
}

func init() {
	registerCollector("fsa", defaultEnabled, NewFSACollector) //defaultDisabled
}

//NewFSACollector -
func NewFSACollector() (Collector, error) {
	return &fsaCollector{
		clusterFSAAdsCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_ads_count"),
			"FSA: Number of alternate data streams.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		clusterFSADirCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_dir_count"),
			"FSA: Number of directories.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		clusterFSAFileCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_file_count"),
			"FSA: Number of files.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		// clusterFSAHasSubdirs: prometheus.NewDesc(
		// 	prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_has_subdirs"),
		// 	"FSA: Defines if directory has subdirectories.",
		// 	[]string{"name", "lin", "parent"}, ConstLabels,
		// ),
		clusterFSALogSizeSum: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_log_size_sum"),
			"FSA: Logical size directory in bytes.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		clusterFSALogSizeSumOverflow: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_log_size_sum_overflow"),
			"FSA: Logical size sum of overflow in bytes.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		clusterFSAOtherCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_other_count"),
			"FSA: Other count.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
		clusterFSAPhysSizeSum: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "fsa_phys_size_sum"),
			"FSA: Physical size directory in bytes.",
			[]string{"name", "lin", "parent"}, ConstLabels,
		),
	}, nil
}

// "Logical size directory in bytes."

func (c *fsaCollector) Update(ch chan<- prometheus.Metric) error {

	// Get FSA result ID
	rsid, b := cache.CacheInstance.Get("rsid")
	if !b {
		//Get Max result-set-id
		resp, err := isiclient.GetFSAResults(IsiCluster.Client)

		if err != nil {
			log.Warnf("Unable to get FSA results. %s", err)
			return err
		}
		rsid = resp.Results[0].ID
		cache.CacheInstance.Set("rsid", rsid, time.Hour)
		log.Infof("ID(rsid): %d", rsid)
	}
	log.Debugf("ID(rsid): %d", rsid)

	resp, err := isiclient.GetFSADirectories(IsiCluster.Client, rsid.(int))
	if err != nil {
		log.Warnf("Unable to get FSA Directories. %s", err)
		return err
	}

	// Total usage
	totalUsage := resp.TotalUsage
	ch <- prometheus.MustNewConstMetric(c.clusterFSAAdsCount, prometheus.GaugeValue, totalUsage.AdsCnt, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSADirCount, prometheus.GaugeValue, totalUsage.DirCnt, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSAFileCount, prometheus.GaugeValue, totalUsage.FileCnt, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSALogSizeSum, prometheus.GaugeValue, totalUsage.LogSizeSum, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSALogSizeSumOverflow, prometheus.GaugeValue, totalUsage.LogSizeSumOverflow, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSAOtherCount, prometheus.GaugeValue, totalUsage.OtherCnt, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))
	ch <- prometheus.MustNewConstMetric(c.clusterFSAPhysSizeSum, prometheus.GaugeValue, totalUsage.PhysSizeSum, totalUsage.Name, strconv.Itoa(totalUsage.LIN), strconv.Itoa(totalUsage.Parent))

	// Usage Data
	for _, usageData := range resp.UsageData {
		ch <- prometheus.MustNewConstMetric(c.clusterFSAAdsCount, prometheus.GaugeValue, usageData.AdsCnt, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSADirCount, prometheus.GaugeValue, usageData.DirCnt, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSAFileCount, prometheus.GaugeValue, usageData.FileCnt, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSALogSizeSum, prometheus.GaugeValue, usageData.LogSizeSum, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSALogSizeSumOverflow, prometheus.GaugeValue, usageData.LogSizeSumOverflow, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSAOtherCount, prometheus.GaugeValue, usageData.OtherCnt, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
		ch <- prometheus.MustNewConstMetric(c.clusterFSAPhysSizeSum, prometheus.GaugeValue, usageData.PhysSizeSum, usageData.Name, strconv.Itoa(usageData.LIN), strconv.Itoa(usageData.Parent))
	}
	return err
}
