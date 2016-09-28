package main

import (
	"github.com/aksentyev/hubble/hubble"
	"github.com/aksentyev/hubble/backend/consul"

	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"

	"./exporter"
	"./util"

    "flag"
    "errors"
    "fmt"

	"net/http"
    // _ "net/http/pprof"
)

// landingPage contains the HTML served at '/'.
// TODO: Make this nicer and more informative.
var landingPage = []byte(`<html>
<head><title>Postgres exporter</title></head>
<body>
<h1>Postgres exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)


var d *hubble.Dispatcher

var (
    consulURL = flag.String(
        "consul.url", "consul.service.consul:8500",
        "Consul url",
    )
    consulDC = flag.String(
        "consul.dc", "staging",
        "Consul datacenter",
    )
	listenAddress = flag.String(
		"listen", ":9113",
		"Address to listen on for web interface and telemetry.",
	)
	metricPath = flag.String(
		"web.telemetry-path", "/metrics",
		"Path under which to expose exporter.",
	)
	queriesPath = flag.String(
		"queries-path", "./queries.yaml",
		"Path to custom queries to run.",
	)
    updateInterval = flag.Int(
        "update-interval", 120,
        "Update interval in seconds",
    )
    scrapeInterval = flag.Int(
        "scrape-interval", 60,
        "Scrape interval in seconds",
    )
)

func setup() {
    config := consul.DefaultConfig()
    config.Address = *consulURL
    config.Datacenter = *consulDC

	client, _ := consul.New(config)

	kv := consul.NewKV(client)
	h := hubble.New(client, kv, "goro")

    filterCB := func(list []*hubble.Service) []*hubble.Service {
        var servicesForMonitoring []*hubble.Service
        for _, svc := range list {
            if util.IncludesStr(svc.Tags, "goro") {
                servicesForMonitoring = append(servicesForMonitoring, svc)
            }
        }
        return servicesForMonitoring
    }

	cb := func() (list []*hubble.ServiceAtomic, err error) {
        defer func() {
            if r := recover(); r != nil {
                err = errors.New(fmt.Sprintf("Unable to get services from consul: %v", r))
                list = []*hubble.ServiceAtomic{}
                log.Errorln(err)
            }
        }()

		for _, svc := range h.Services(filterCB){
			for _, el := range svc.MakeAtomic(nil) {
				list = append(list, el)
			}
		}
		return list, err
	}

	d = h.NewDispatcher(*updateInterval, cb)
}

func main(){
    // Profiler
    // go func() {
	//     http.ListenAndServe("localhost:6060", nil)
    // }()
	flag.Parse()
	setup()

    pgMetricsParsed := exporter.AddFromFile(*queriesPath)

	go func() {
		for svc := range d.ToRegister {
            if len(svc.ExporterOptions) > 1 {
                config := exporter.Config{
                    DSN:             util.PgConnURL(svc),
                    Labels:          svc.ExtraLabels,
                    ExporterOptions: svc.ExporterOptions,
                    CacheTTL:        *scrapeInterval,
                    PgMetrics:       pgMetricsParsed,
                }
                exp, err := exporter.CreateAndRegister(&config)
                if err == nil {
                    d.Register(svc, exp)
                    log.Infof("Registered %v %v", svc.Name, svc.Address)
                } else {
                    log.Infof("Register was failed for service %v %v %v", svc.Name, svc.Address, err)
                    exp.Close()
                }
            }
		}
	}()

    go func() {
        for m := range d.ToUnregister {
            for h, svc := range m {
                exporter := d.Exporters[h].(*exporter.PostgresExporter)
                err := exporter.Close()
                if err != nil {
                    log.Warnf("Unregister() for %v %v returned %v:", svc.Name, svc.Address, err)
                } else {
                    log.Infof("Unregister service %v %v", svc.Name, svc.Address)
                }
                d.UnregisterWithHash(h)
            }
		}
    }()

	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infof("Starting Server: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
