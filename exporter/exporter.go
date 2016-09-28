package exporter

import (
    "database/sql"
    _ "github.com/lib/pq"

    // "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"

    "github.com/aksentyev/hubble/exportertools"
)

// Exporter collects Postgres metrics. It implements prometheus.Collector.
type PostgresExporter struct {
    *exportertools.BaseExporter
    Config   *Config
    db       *sql.DB
}

// NewExporter returns a new PostgreSQL exporter for the provided DSN.
func CreateAndRegister(config *Config) (*PostgresExporter, error) {
    exp := PostgresExporter{
        Config: config,
        BaseExporter: exportertools.NewBaseExporter("postgres", config.CacheTTL, config.Labels),
    }
    err := exportertools.Register(&exp)
    if err != nil {
        return &exp, err
    }
    return &exp, nil
}

func (e *PostgresExporter) Setup() error {
    db, err := sql.Open("postgres", e.Config.DSN)
    if err != nil {
        log.Infoln("Error opening connection to database:", err)
        return err
    }
    e.db = db
    e.AddCollector(NewCollector(db, e.Config))
    return nil
}

func (e *PostgresExporter) Close() (err error) {
    defer close(e.Control)

    err = exportertools.Unregister(e)
    if e.db != nil {
        err = e.db.Close()
        log.Infoln("db closed")
    }

    e.Control<- true
    log.Debugf("Stop processing metric for %v", e.Labels)

    return err
}
