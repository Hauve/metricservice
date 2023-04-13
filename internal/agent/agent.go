package agent

import (
	"flag"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/handlers"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type MyAgent struct {
	storage handlers.Storage

	address        string
	reportInterval int
	pollInterval   int
}

func New() *MyAgent {
	address := flag.String("a", "localhost:8080", "address")
	reportInterval := flag.Int("r", 10, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")

	flag.Parse()

	return &MyAgent{
		storage:        storage.NewMemStorage(),
		address:        *address,
		reportInterval: *reportInterval,
		pollInterval:   *pollInterval,
	}
}

func (ag *MyAgent) Run() {
	dur := fmt.Sprintf("%ds", ag.pollInterval)
	duration, err := time.ParseDuration(dur)
	if err != nil {
		log.Printf("%e", err)
	}
	go ag.sendMetrics()
	for {
		time.Sleep(duration * time.Second)
		ag.collectMetrics()
	}
}

func (ag *MyAgent) collectMetrics() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	ag.storage.SetGauge("Alloc", float64(stats.Alloc))
	ag.storage.SetGauge("BuckHashSys", float64(stats.BuckHashSys))
	ag.storage.SetGauge("Frees", float64(stats.Frees))
	ag.storage.SetGauge("GCCPUFraction", stats.GCCPUFraction)
	ag.storage.SetGauge("GCSys", float64(stats.GCSys))
	ag.storage.SetGauge("HeapAlloc", float64(stats.HeapAlloc))
	ag.storage.SetGauge("HeapIdle", float64(stats.HeapIdle))
	ag.storage.SetGauge("HeapInuse", float64(stats.HeapInuse))
	ag.storage.SetGauge("HeapObjects", float64(stats.HeapObjects))
	ag.storage.SetGauge("HeapReleased", float64(stats.HeapReleased))
	ag.storage.SetGauge("HeapSys", float64(stats.HeapSys))
	ag.storage.SetGauge("LastGC", float64(stats.LastGC))
	ag.storage.SetGauge("Lookups", float64(stats.Lookups))
	ag.storage.SetGauge("MCacheInuse", float64(stats.MCacheInuse))
	ag.storage.SetGauge("MCacheSys", float64(stats.MCacheSys))
	ag.storage.SetGauge("MSpanInuse", float64(stats.MSpanInuse))
	ag.storage.SetGauge("MSpanSys", float64(stats.MSpanSys))
	ag.storage.SetGauge("Mallocs", float64(stats.Mallocs))
	ag.storage.SetGauge("NextGC", float64(stats.NextGC))
	ag.storage.SetGauge("NumForcedGC", float64(stats.NumForcedGC))
	ag.storage.SetGauge("NumGC", float64(stats.NumGC))
	ag.storage.SetGauge("OtherSys", float64(stats.OtherSys))
	ag.storage.SetGauge("PauseTotalNs", float64(stats.PauseTotalNs))
	ag.storage.SetGauge("StackInuse", float64(stats.StackInuse))
	ag.storage.SetGauge("StackSys", float64(stats.StackSys))
	ag.storage.SetGauge("Sys", float64(stats.Sys))
	ag.storage.SetGauge("OtherSys", float64(stats.TotalAlloc))
	ag.storage.SetGauge("RandomValue", float64(stats.TotalAlloc))

	ag.storage.SetGauge("RandomValue", rand.Float64())

	pollCount, ok := ag.storage.GetCounter("PollCount")
	if !ok {
		ag.storage.SetCounter("PoolCount", 1)
		return
	}
	ag.storage.SetCounter("PollCount", pollCount)

}

func (ag *MyAgent) sendMetrics() {
	dur := fmt.Sprintf("%ds", ag.reportInterval)
	duration, err := time.ParseDuration(dur)
	if err != nil {
		log.Printf("%e", err)
	}

	for {
		time.Sleep(duration)
		gaugeNames := ag.storage.GetGaugeKeys()
		for _, name := range *gaugeNames {
			value, _ := ag.storage.GetGauge(name)
			val := fmt.Sprintf("%f", value)
			ag.sendMetric(name, val, storage.Gauge)
		}

		counterNames := ag.storage.GetCounterKeys()
		for _, name := range *counterNames {
			value, _ := ag.storage.GetCounter(name)
			val := fmt.Sprintf("%d", value)
			ag.sendMetric(name, val, storage.Gauge)
		}
	}
}

func (ag *MyAgent) sendMetric(name string, value string, mt storage.MetricType) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", ag.address, mt, name, value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Printf("An Error Occured %v\n", err)
		return
	}
	req.Header.Add("Content-Length", `0`)
	req.Header.Add("Content-Type", `text/plain`)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("An Error Occured %v\n", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Printf("An Error Occured %v\n", err)
		return
	}
}
