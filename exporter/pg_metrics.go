package exporter

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/prometheus/common/log"
)

type PgMetric struct {
    query string
    specs []*PgMetricSpecs
}

type PgMetricSpecs struct {
    name, desc, mtype string
}

// AddFromFile reads yaml with metric parameter and serialize it
func AddFromFile(queriesPath string) (metrics []*PgMetric) {
    queries := loadFile(queriesPath)
    for prefix, specs := range queries {
        metric := PgMetric{}
        for key, value := range specs.(map[interface{}]interface{}) {
            switch key.(string) {
            case "query":
                query := value.(string)
                metric.query = query

            case "metrics":
                for _, item := range value.([]interface{}) {
                    metricPropsInterface := item.(map[interface{}]interface{})
                    metricProps := PgMetricSpecs{}

                    for name, attrs := range metricPropsInterface {
                        metricProps.name = fmt.Sprintf("%v_%v",prefix, name.(string))

                        for key, val := range attrs.(map[interface{}]interface{}) {
                            switch key.(string) {
                            case "type":
                                metricProps.mtype = val.(string)
                            case "description":
                                metricProps.desc = val.(string)
                            }
                        }
                    }
                    metric.specs = append(metric.specs, &metricProps)
                }
            }
        }
        metrics = append(metrics, &metric)
    }
    log.Infof("Metrics parsed:\n%v", spew.Sdump(metrics))
    return metrics
}

func loadFile(path string) map[string]interface{} {
    var extra map[string]interface{}

    content, err := ioutil.ReadFile(path)
    if err != nil {
        log.Errorf("File read error: %v", err)
    }

    err = yaml.Unmarshal(content, &extra)
    if err != nil {
        log.Errorf("File parse error: %v", err)
    }

    return extra
}
