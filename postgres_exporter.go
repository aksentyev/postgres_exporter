package main

import (
    "github.com/aksentyev/hubble/hubble"
    "github.com/aksentyev/hubble/backend/consul"
    "github.com/aksentyev/hubble/exportertools"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"

    "github.com/aksentyev/postgres_exporter/exporter"
    "github.com/aksentyev/postgres_exporter/util"

    "flag"
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
    consulTag = flag.String(
        "consul.tag", "postgres",
        "Look for services that have the tag specified.",
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
    showVersion = flag.Bool(
        "version", false,
        "Show versions and exit",
    )
)

func setup() {
    config := consul.DefaultConfig()
    config.Address = *consulURL
    config.Datacenter = *consulDC

    client, err := consul.New(config)
    if err != nil {
        panic(err)
    }

    kv := consul.NewKV(client)
    h := hubble.New(client, kv, *consulTag)

    filterCB := func(list []*hubble.Service) []*hubble.Service {
        var servicesForMonitoring []*hubble.Service
        for _, svc := range list {
            if util.IncludesStr(svc.Tags, *consulTag) {
                servicesForMonitoring = append(servicesForMonitoring, svc)
            }
        }
        return servicesForMonitoring
    }

    cb := func() (list []*hubble.ServiceAtomic, err error) {
        services, err := h.Services(filterCB)
        if err != nil {
            return list, err
        }
        for _, svc := range services {
            for _, el := range svc.MakeAtomic(nil) {
                list = append(list, el)
            }
        }
        return list, err
    }

    d = hubble.NewDispatcher(*updateInterval)
    go d.Run(cb)
}

func printVersions(){
    fmt.Printf("exporter: %v\n", exporter.VERSION)
    fmt.Printf("hubble: %v\n", hubble.VERSION)
    fmt.Printf("exportertools: %v\n", exportertools.VERSION)
    fmt.Printf("consul backend: %v\n", consul.VERSION)
}

func listenAndRegister() {
    pgMetricsParsed := exporter.AddFromFile(*queriesPath)

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
                log.Warnf("Register was failed for service %v %v %v", svc.Name, svc.Address, err)
                exp.Close()
            }
        }
    }
}

func listenAndUnregister() {
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
}

func main(){
    // Profiler
    // go func() {
    //     http.ListenAndServe("localhost:6060", nil)
    // }()
    flag.Parse()

    if *showVersion {
        printVersions()
        return
    }

    setup()
    go listenAndRegister()
    go listenAndUnregister()

    http.Handle(*metricPath, prometheus.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write(landingPage)
    })
    log.Infof("Starting Server: %s", *listenAddress)
    log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
