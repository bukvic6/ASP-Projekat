package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	currentCount = 0

	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of http hits.",
		},
	)
	createConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_config_hit_total",
			Help: "Total number od all create config hits",
		})

	getAllHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_hit_total",
			Help: "Total number of get all hits",
		})
	getConfigVersionsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_versions_hit_total",
			Help: "Total number of get config versions hits",
		})
	getConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_Config_hit_total",
			Help: "Total number of get config hits",
		})
	addConfigVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "add_Config_version_hit_total",
			Help: "Total number of get config version hits",
		})
	delConfigVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_version_hit_total",
			Help: "Total number of del config version hits",
		})
	createGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "Create_group_hit_total",
			Help: "Total number of Create group hits",
		})
	getAllGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_group_hit_total",
			Help: "Total number of get all group hits",
		})
	addGroupVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "add_group_version_hit_total",
			Help: "Total number of add group version hits",
		})
	getConfigGroupVersionsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_versions_hit_total",
			Help: "Total number of get group versions hits",
		})
	getGroupVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_version_hit_total",
			Help: "Total number of get group version hits",
		})
	delgroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_group_hit_total",
			Help: "Total number of del group hits",
		})
	addConfigToGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "add_config_to_group_hit_total",
			Help: "Total number of add config to group hits",
		})
	metricsList = []prometheus.Collector{
		createConfigHits, getAllHits, getConfigVersionsHits, getConfigHits,
		addConfigVersionHits, delConfigVersionHits, createGroupHits, getAllGroupHits,
		addGroupVersionHits, getConfigGroupVersionsHits, getGroupVersionHits, delgroupHits,
		addConfigToGroupHits,
	}
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	prometheusRegistry.MustRegister(metricsList...)
}
func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}
func counteAddConfigToGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		addConfigToGroupHits.Inc()
		f(w, r) // original function call
	}
}
func counteDelgroupHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delgroupHits.Inc()
		f(w, r) // original function call
	}
}
func counteGetGroupVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupVersionHits.Inc()
		f(w, r) // original function call
	}
}
func counteGetConfigGroupVersions(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigGroupVersionsHits.Inc()
		f(w, r) // original function call
	}
}
func counteAddGroupVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		addGroupVersionHits.Inc()
		f(w, r) // original function call
	}
}
func countegetAllGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllGroupHits.Inc()
		f(w, r) // original function call
	}
}
func counteCreateGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createGroupHits.Inc()
		f(w, r) // original function call
	}
}
func countdelConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigVersionHits.Inc()
		f(w, r) // original function call
	}
}
func countAddConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		addConfigVersionHits.Inc()
		f(w, r) // original function call
	}
}
func countGetConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigHits.Inc()
		f(w, r) // original function call
	}
}
func countConfigVersions(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigVersionsHits.Inc()
		f(w, r) // original function call
	}
}
func countCreateConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createConfigHits.Inc()
		f(w, r) // original function call
	}
}
func countGetAll(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllHits.Inc()
		f(w, r) // original function call
	}
}
