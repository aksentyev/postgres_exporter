package exporter

import (
    "github.com/aksentyev/hubble/exportertools"
)

// statistics for database name provided
type DatabaseStat struct {
    Name                                          string
    Commit, Rollback                              int64
    Read, Hit                                     int64
    Returned, Fetched, Inserted, Updated, Deleted int64
    Conflicts, TempFiles, TempBytes, Deadlocks    int64
    ReadTime, WriteTime                           float64
}

func formatDatabaseStats(labels map[string]string, s *DatabaseStat) []*exportertools.Metric {
    metrics := []*exportertools.Metric{
        {
            Name:        "commits",
            Type:        exportertools.Gauge,
            Value:       s.Commit,
            Description: "commits",
            Labels:      labels,
        },
        {
            Name:        "rollbacks",
            Type:        exportertools.Gauge,
            Value:       s.Rollback,
            Description: "rollbacks",
            Labels:      labels,
        },
        {
            Name:        "reads",
            Type:        exportertools.Gauge,
            Value:       s.Read,
            Description: "reads",
            Labels:      labels,
        },
        {
            Name:        "hits",
            Type:        exportertools.Gauge,
            Value:       s.Hit,
            Description: "hits",
            Labels:      labels,
        },
        {
            Name:        "returns",
            Type:        exportertools.Gauge,
            Value:       s.Commit,
            Description: "returns",
            Labels:      labels,
        },
        {
            Name:        "fetches",
            Type:        exportertools.Gauge,
            Value:       s.Fetched,
            Description: "fetches",
            Labels:      labels,
        },
        {
            Name:        "inserts",
            Type:        exportertools.Gauge,
            Value:       s.Inserted,
            Description: "inserts",
            Labels:      labels,
        },
        {
            Name:        "updates",
            Type:        exportertools.Gauge,
            Value:       s.Updated,
            Description: "updates",
            Labels:      labels,
        },
        {
            Name:        "deletes",
            Type:        exportertools.Gauge,
            Value:       s.Deleted,
            Description: "deletes",
            Labels:      labels,
        },
        {
            Name:        "conflicts",
            Type:        exportertools.Gauge,
            Value:       s.Conflicts,
            Description: "conflicts",
            Labels:      labels,
        },
        {
            Name:        "temp_files",
            Type:        exportertools.Gauge,
            Value:       s.TempFiles,
            Description: "temp_files",
            Labels:      labels,
        },
        {
            Name:        "temp_bytes",
            Type:        exportertools.Gauge,
            Value:       s.TempBytes,
            Description: "temp_bytes",
            Labels:      labels,
        },
        {
            Name:        "deadlocks",
            Type:        exportertools.Gauge,
            Value:       s.Commit,
            Description: "deadlocks",
            Labels:      labels,
        },
    }

    return metrics
}
