package exporter

import (
    "database/sql"
    _ "github.com/lib/pq"
    "strconv"

    "github.com/aksentyev/hubble/exportertools"
)

type Collector struct {
    db        *sql.DB
    *Config
}

func NewCollector(db *sql.DB, config *Config) *Collector {
    return &Collector{db, config}
}

// Collecting metrics
func (c *Collector) Collect() ([]*exportertools.Metric, error) {
    stat := new(DatabaseStat)
    err := c.collectDatabaseStats(c.ExporterOptions["db"], stat) //TODO:
    if err != nil {
        return make([]*exportertools.Metric, 0), err
    }
    dbMetrics := formatDatabaseStats(c.Labels, stat)
    pgMetrics := c.collectPgStats()
    return append(dbMetrics, pgMetrics...), nil
}

func (c *Collector) collectDatabaseStats(dbName string, s *DatabaseStat) error {
    var ignore interface{}
    return c.db.QueryRow(`SELECT * FROM pg_stat_database WHERE datname=$1`, dbName).Scan(
        &ignore, &s.Name, &ignore,
        &s.Commit, &s.Rollback,
        &s.Read, &s.Hit,
        &s.Returned, &s.Fetched, &s.Inserted, &s.Updated, &s.Deleted,
        &s.Conflicts, &s.TempFiles, &s.TempBytes, &s.Deadlocks,
        &s.ReadTime, &s.WriteTime,
        &ignore,
    )
}

func (c *Collector) collectPgStats() (collectedData []*exportertools.Metric) {
    for _, m := range c.PgMetrics {
        value, err := c.getValue(m.query)
        if err != nil {
            continue
        }

        em := exportertools.Metric{
                Name:        m.name,
                Type:        exportertools.StringToType(m.mtype),
                Value:       value,
                Description: m.desc,
                Labels:      c.Labels,
            }

        collectedData = append(collectedData, &em)
    }
    return collectedData
}

func (c *Collector) getValue(query string) (int64, error) {
    var data string
    err := c.db.QueryRow(query).Scan(&data)
    if err != nil {
        return -1, err
    }
    v, _ := strconv.ParseInt(data, 10, 64)
    return v, nil
}
