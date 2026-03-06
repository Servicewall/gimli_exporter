package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/process"
)

type gimliCollector struct {
	cpuUsage    *prometheus.Desc
	memoryUsage *prometheus.Desc
	netConns    *prometheus.Desc
}

func newGimliCollector() *gimliCollector {
	return &gimliCollector{
		cpuUsage: prometheus.NewDesc(
			"gimli_cpu_usage_percent",
			"CPU usage percentage of the gimli process",
			[]string{"pid"}, nil,
		),
		memoryUsage: prometheus.NewDesc(
			"gimli_memory_bytes",
			"Memory usage (RSS) in bytes of the gimli process",
			[]string{"pid"}, nil,
		),
		netConns: prometheus.NewDesc(
			"gimli_net_connections_total",
			"Total number of network connections of the gimli process",
			[]string{"pid"}, nil,
		),
	}
}

func (c *gimliCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuUsage
	ch <- c.memoryUsage
	ch <- c.netConns
}

func (c *gimliCollector) Collect(ch chan<- prometheus.Metric) {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error fetching processes: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if name == "gimli" || name == "gimli.exe" {
			wg.Add(1)
			go func(p *process.Process) {
				defer wg.Done()
				pidStr := fmt.Sprintf("%d", p.Pid)

				// CPU Percent: Differentiate between OS via getCPUPercent
				cpu, err := getCPUPercent(p)
				if err == nil {
					ch <- prometheus.MustNewConstMetric(c.cpuUsage, prometheus.GaugeValue, cpu, pidStr)
				}

				// Memory Info (RSS)
				mem, err := p.MemoryInfo()
				if err == nil {
					ch <- prometheus.MustNewConstMetric(c.memoryUsage, prometheus.GaugeValue, float64(mem.RSS), pidStr)
				}

				// Network Connections
				conns, err := p.Connections()
				if err == nil {
					ch <- prometheus.MustNewConstMetric(c.netConns, prometheus.GaugeValue, float64(len(conns)), pidStr)
				}
			}(p)
		}
	}
	wg.Wait()
}

func main() {
	collector := newGimliCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Gimli Exporter</title></head>
			<body>
			<h1>Gimli Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	port := ":9101"
	log.Printf("Starting gimli_exporter on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
