package exporter

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "fmt"
    "github.com/prometheus/common/log"
)

type PgMetric struct {
    query, name, desc, mtype string
}

// AddFromFile reads yaml with metric parameter and serialize it
func AddFromFile(queriesPath string) (metrics []*PgMetric) {
    var extra map[string]interface{}

    content, err := ioutil.ReadFile(queriesPath)
    if err != nil {
        log.Errorf("File read error: %v", err)
    }

    err = yaml.Unmarshal(content, &extra)
    if err != nil {
        log.Errorf("File parse error: %v", err)
    }

    for prefix, specs := range extra {
        metricItem := PgMetric{}
        for key, value := range specs.(map[interface{}]interface{}) {
            switch key.(string) {
            case "query":
                query := value.(string)
                metricItem.query = query

            case "metrics":
                for _, c := range value.([]interface{}) {
                    metricProps := c.(map[interface{}]interface{})

                    for name, attrs := range metricProps {
                        metricItem.name = fmt.Sprintf("%v_%v",prefix, name.(string))

                        for key, val := range attrs.(map[interface{}]interface{}) {
                            switch key.(string) {
                            case "type":
                                metricItem.mtype = val.(string)
                            case "description":
                                metricItem.desc = val.(string)
                            }
                        }
                    }
                }
            }
        }
        metrics = append(metrics, &metricItem)
    }
    return metrics
}
