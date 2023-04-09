package main

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var AgentStorage = storage.NewMemStorage()

func collectMetrics() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	AgentStorage.SetGauge("Alloc", float64(stats.Alloc))
	AgentStorage.SetGauge("BuckHashSys", float64(stats.BuckHashSys))
	AgentStorage.SetGauge("Frees", float64(stats.Frees))
	AgentStorage.SetGauge("GCCPUFraction", stats.GCCPUFraction)
	AgentStorage.SetGauge("GCSys", float64(stats.GCSys))
	AgentStorage.SetGauge("HeapAlloc", float64(stats.HeapAlloc))
	AgentStorage.SetGauge("HeapIdle", float64(stats.HeapIdle))
	AgentStorage.SetGauge("HeapInuse", float64(stats.HeapInuse))
	AgentStorage.SetGauge("HeapObjects", float64(stats.HeapObjects))
	AgentStorage.SetGauge("HeapReleased", float64(stats.HeapReleased))
	AgentStorage.SetGauge("HeapSys", float64(stats.HeapSys))
	AgentStorage.SetGauge("LastGC", float64(stats.LastGC))
	AgentStorage.SetGauge("Lookups", float64(stats.Lookups))
	AgentStorage.SetGauge("MCacheInuse", float64(stats.MCacheInuse))
	AgentStorage.SetGauge("MCacheSys", float64(stats.MCacheSys))
	AgentStorage.SetGauge("MSpanInuse", float64(stats.MSpanInuse))
	AgentStorage.SetGauge("MSpanSys", float64(stats.MSpanSys))
	AgentStorage.SetGauge("Mallocs", float64(stats.Mallocs))
	AgentStorage.SetGauge("NextGC", float64(stats.NextGC))
	AgentStorage.SetGauge("NumForcedGC", float64(stats.NumForcedGC))
	AgentStorage.SetGauge("NumGC", float64(stats.NumGC))
	AgentStorage.SetGauge("OtherSys", float64(stats.OtherSys))
	AgentStorage.SetGauge("PauseTotalNs", float64(stats.PauseTotalNs))
	AgentStorage.SetGauge("StackInuse", float64(stats.StackInuse))
	AgentStorage.SetGauge("StackSys", float64(stats.StackSys))
	AgentStorage.SetGauge("Sys", float64(stats.Sys))
	AgentStorage.SetGauge("OtherSys", float64(stats.TotalAlloc))
	AgentStorage.SetGauge("RandomValue", float64(stats.TotalAlloc))

	AgentStorage.SetGauge("RandomValue", rand.Float64())

	pollCount, ok := AgentStorage.GetCounter("PollCount")
	if !ok {
		AgentStorage.SetCounter("PoolCount", 1)
		return
	}
	AgentStorage.SetCounter("PollCount", pollCount)

}

func sendMetrics() {

	gaugeNames := AgentStorage.GetGaugeKeys()
	for _, name := range *gaugeNames {
		value, _ := AgentStorage.GetGauge(name)
		val := fmt.Sprintf("%f", value)
		go sendMetric(name, val, storage.Gauge)
	}

	counterNames := AgentStorage.GetCounterKeys()
	for _, name := range *counterNames {
		value, _ := AgentStorage.GetCounter(name)
		val := fmt.Sprintf("%d", value)
		go sendMetric(name, val, storage.Gauge)
	}
}

func sendMetric(name string, value string, mt storage.MetricType) {
	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", mt, name, value)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("Content-Length", `0`)
	req.Header.Add("Content-Type", `text/plain`)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("An Error Occured %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("An Error Occured %v\n", err)
			return
		}
	}(resp.Body)
}

func main() {

	for {
		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)
			collectMetrics()
		}
		sendMetrics()
	}
}
