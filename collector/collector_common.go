/*
Copyright 2018 Adobe
All Rights Reserved.

NOTICE: Adobe permits you to use, modify, and distribute this file in
accordance with the terms of the Adobe license agreement accompanying
it. If you have received this file from a source other than Adobe,
then your use, modification, or distribution of it requires the prior
written permission of Adobe.
*/

package collector

import (
	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"

	"github.com/hpanike/goisilon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	//IsiCluster is the structure that holds all the information need to establish the connection.
	IsiCluster IsilonCluster
	//ConstLabels are constant labels that every metric will have.  This includes the label cluster.
	ConstLabels prometheus.Labels

	statsEngineCallDuration *prometheus.Desc
	statsEngineCallFailure  *prometheus.Desc
	exporterDurationDesc    *prometheus.Desc
	scrapeSuccessDesc       *prometheus.Desc
	scrapeDurationDesc      *prometheus.Desc
)

//IsilonCluster struct contains all the connection info and an instanciated client connection to the cluster.
type IsilonCluster struct {
	FQDN        string
	Name        string
	Port        string
	Username    string
	Site        string
	PasswordEnv string
	QuotaOnly   bool
	Quotas      Quotas
	FSAOnly     bool
	Client      *goisilon.Client
}

//Quotas struct contains information for to quota only collections
type Quotas struct {
	Count  int64
	Errors int64
	Err    error
	Retry  int64
}

//SetClusterConfigName will get the name from the isi config and set it as IsilonClusterConfigName inside IsiCluster.
func SetClusterConfigName() error {
	clusterName, err := isiclient.GetClusterName(IsiCluster.Client)
	if err != nil {
		log.Warnf("Unabled to obtain cluster name from isi config.")
		return err
	}
	IsiCluster.Name = clusterName

	err = CreateConstLabels()
	if err != nil {
		log.Warnf("Unable to create const labels.")
	}
	return nil
}

//GetClusterConnector calls the isiclient and creates a new isilon cluster connector.
func GetClusterConnector() error {
	con, err := isiclient.NewIsilonClient(IsiCluster.FQDN, IsiCluster.Port, IsiCluster.Username, IsiCluster.PasswordEnv)
	if err != nil {
		log.Warn("Unabled to create connection to the Isilon cluster.")
		return err
	}
	IsiCluster.Client = con
	return nil
}

//CreateConstLabels will create an array of labels that are constant to all metrics.
func CreateConstLabels() error {
	//Only create a const label for site if a site has been specified.
	if IsiCluster.Site != "" {
		ConstLabels = prometheus.Labels{"cluster": IsiCluster.Name, "site": IsiCluster.Site}
	} else {
		ConstLabels = prometheus.Labels{"cluster": IsiCluster.Name}
	}
	log.Debugf("ConstLables are %v", ConstLabels)
	return nil
}

//GetNumQuotas retrieve the number of quotas the system should have.
func GetNumQuotas() error {
	summary, err := isiclient.GetQuotaSummary(IsiCluster.Client)
	if err != nil {
		log.Warn("Unabled to update quota summary information.")
		return err
	}
	IsiCluster.Quotas.Count = int64(summary.Count)

	return nil
}
