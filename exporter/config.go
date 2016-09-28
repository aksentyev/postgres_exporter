package exporter

type Config struct {
    DSN             string
    Labels          map[string]string
    ExporterOptions map[string]string
    CacheTTL        int
    PgMetrics       []*PgMetric
}
