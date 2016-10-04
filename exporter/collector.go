package exporter

import (
    "database/sql"
    _ "github.com/lib/pq"

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
    dbData := formatDatabaseStats(c.Labels, stat)
    pgData := c.collectPgStats()
    return append(dbData, pgData...), nil
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
        values := c.db.QueryRow(m.query)

        fields := tableFields(m.specs)
        parsedData := make([]interface{}, len(fields))
        scanArgs := make([]interface{}, len(fields))
        for i := range parsedData {
            scanArgs[i] = &parsedData[i]
        }

        err := values.Scan(scanArgs...)
        if err != nil {
            continue
        }

        for idx, val := range parsedData {

            em := exportertools.Metric{
                    Name:        m.specs[idx].name,
                    Type:        exportertools.StringToType(m.specs[idx].mtype),
                    Value:       val,
                    Description: m.specs[idx].desc,
                    Labels:      c.Labels,
                }

            collectedData = append(collectedData, &em)
        }

    }
    return collectedData
}

func tableFields(specs []*PgMetricSpecs) (list []string) {
    for _, s := range specs {
        list = append(list, s.name)
    }
    return list
}
