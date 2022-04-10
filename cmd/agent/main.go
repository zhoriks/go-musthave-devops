package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var memory runtime.MemStats
	pollCount := 0
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	tickerMetric := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	signalChanel := make(chan os.Signal, 1)
	endpoint := "http://127.0.0.1:8080/update"

	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	type Metric struct {
		Name  string
		Type  string
		Value string
	}
	var metrics [29]Metric

	for {
		select {
		case <-tickerMetric.C:
			rand.Seed(time.Now().UnixNano())
			randomValue := rand.Int()
			runtime.ReadMemStats(&memory)
			pollCount++

			metrics[0] = Metric{Name: "Alloc", Type: "gauge", Value: strconv.FormatUint(memory.Alloc, 10)}
			metrics[1] = Metric{Name: "BuckHashSys", Type: "gauge", Value: strconv.FormatUint(memory.BuckHashSys, 10)}
			metrics[2] = Metric{Name: "Frees", Type: "gauge", Value: strconv.FormatUint(memory.Frees, 10)}
			metrics[3] = Metric{Name: "GCCPUFraction", Type: "gauge", Value: strconv.FormatFloat(memory.GCCPUFraction, 'f', 6, 64)}
			metrics[4] = Metric{Name: "GCSys", Type: "gauge", Value: strconv.FormatUint(memory.GCSys, 10)}
			metrics[5] = Metric{Name: "HeapAlloc", Type: "gauge", Value: strconv.FormatUint(memory.HeapAlloc, 10)}
			metrics[6] = Metric{Name: "HeapIdle", Type: "gauge", Value: strconv.FormatUint(memory.HeapIdle, 10)}
			metrics[7] = Metric{Name: "HeapInuse", Type: "gauge", Value: strconv.FormatUint(memory.HeapInuse, 10)}
			metrics[8] = Metric{Name: "HeapObjects", Type: "gauge", Value: strconv.FormatUint(memory.HeapObjects, 10)}
			metrics[9] = Metric{Name: "HeapReleased", Type: "gauge", Value: strconv.FormatUint(memory.HeapReleased, 10)}
			metrics[10] = Metric{Name: "HeapSys", Type: "gauge", Value: strconv.FormatUint(memory.HeapSys, 10)}
			metrics[11] = Metric{Name: "LastGC", Type: "gauge", Value: strconv.FormatUint(memory.LastGC, 10)}
			metrics[12] = Metric{Name: "Lookups", Type: "gauge", Value: strconv.FormatUint(memory.Lookups, 10)}
			metrics[13] = Metric{Name: "MCacheInuse", Type: "gauge", Value: strconv.FormatUint(memory.MCacheInuse, 10)}
			metrics[14] = Metric{Name: "MCacheSys", Type: "gauge", Value: strconv.FormatUint(memory.MCacheSys, 10)}
			metrics[15] = Metric{Name: "MSpanInuse", Type: "gauge", Value: strconv.FormatUint(memory.MSpanInuse, 10)}
			metrics[16] = Metric{Name: "MSpanSys", Type: "gauge", Value: strconv.FormatUint(memory.MSpanSys, 10)}
			metrics[17] = Metric{Name: "Mallocs", Type: "gauge", Value: strconv.FormatUint(memory.Mallocs, 10)}
			metrics[18] = Metric{Name: "NextGC", Type: "gauge", Value: strconv.FormatUint(memory.NextGC, 10)}
			metrics[19] = Metric{Name: "NumForcedGC", Type: "gauge", Value: strconv.FormatUint(uint64(memory.NumForcedGC), 10)}
			metrics[20] = Metric{Name: "NumGC", Type: "gauge", Value: strconv.FormatUint(uint64(memory.NumGC), 10)}
			metrics[21] = Metric{Name: "OtherSys", Type: "gauge", Value: strconv.FormatUint(memory.OtherSys, 10)}
			metrics[22] = Metric{Name: "PauseTotalNs", Type: "gauge", Value: strconv.FormatUint(memory.PauseTotalNs, 10)}
			metrics[23] = Metric{Name: "StackInuse", Type: "gauge", Value: strconv.FormatUint(memory.StackInuse, 10)}
			metrics[24] = Metric{Name: "StackSysStackSys", Type: "gauge", Value: strconv.FormatUint(memory.StackSys, 10)}
			metrics[25] = Metric{Name: "Sys", Type: "gauge", Value: strconv.FormatUint(memory.Sys, 10)}
			metrics[26] = Metric{Name: "TotalAlloc", Type: "gauge", Value: strconv.FormatUint(memory.TotalAlloc, 10)}
			metrics[27] = Metric{Name: "PollInterval", Type: "gauge", Value: strconv.FormatUint(uint64(pollCount), 10)}
			metrics[28] = Metric{Name: "RandomValue", Type: "counter", Value: strconv.FormatUint(uint64(randomValue), 10)}

		case <-tickerReport.C:
			client := &http.Client{}
			for _, m := range metrics {
				fmt.Println(m)
				request, err := http.NewRequest(
					http.MethodPost,
					endpoint+"/"+m.Type+"/"+m.Name+"/"+m.Value,
					nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				request.Header.Add("Content-Type", "text/plain")

				response, err := client.Do(request)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer response.Body.Close()
			}
		case <-signalChanel:
			os.Exit(1)
		}

		//s := <-signalChanel
		//switch s {
		//case syscall.SIGINT:
		//	fmt.Println("Signal interrupt triggered.")
		//	exitChan <- 0
		//case syscall.SIGTERM:
		//	fmt.Println("Signal terminate triggered.")
		//	exitChan <- 0
		//case syscall.SIGQUIT:
		//	fmt.Println("Signal quit triggered.")
		//	exitChan <- 0
		//}
		//fmt.Printf("Alloc %s,\n", Alloc.Value)
		//fmt.Printf("BuckHashSys %s,\n", BuckHashSys.Value)
		//fmt.Printf("Frees %s\n", Frees.Value)
		//fmt.Printf("GCCPUFraction %s,\n", GCCPUFraction.Value)
		//fmt.Printf("GCSys %s,\n", GCSys.Value)
		//fmt.Printf("HeapAlloc %s,\n", HeapAlloc.Value)
		//fmt.Printf("HeapIdle %s,\n", HeapIdle.Value)
		//fmt.Printf("HeapInuse %s,\n", HeapInuse.Value)
		//fmt.Printf("HeapObjects %s,\n", HeapObjects.Value)
		//fmt.Printf("HeapReleased %s,\n", HeapReleased.Value)
		//fmt.Printf("HeapSys %s,\n", HeapSys.Value)
		//fmt.Printf("LastGC %s,\n", LastGC.Value)
		//fmt.Printf("Lookups %s,\n", Lookups.Value)
		//fmt.Printf("MCacheInuse %s,\n", MCacheInuse.Value)
		//fmt.Printf("MCacheSys %s,\n", MCacheSys.Value)
		//fmt.Printf("MSpanInuse %s,\n", MSpanInuse.Value)
		//fmt.Printf("MSpanSys %s,\n", MSpanSys.Value)
		//fmt.Printf("Mallocs %s,\n", Mallocs.Value)
		//fmt.Printf("NextGC %s,\n", NextGC.Value)
		//fmt.Printf("NumForcedGC %s,\n", NumForcedGC.Value)
		//fmt.Printf("NumGC %s,\n", NumGC.Value)
		//fmt.Printf("OtherSys %s,\n", OtherSys.Value)
		//fmt.Printf("PauseTotalNs %s,\n", PauseTotalNs.Value)
		//fmt.Printf("StackInuse %s,\n", StackInuse.Value)
		//fmt.Printf("StackSys %s,\n", StackSys.Value)
		//fmt.Printf("Sys %s,\n", Sys.Value)
		//fmt.Printf("TotalAlloc %s,\n", TotalAlloc.Value)
		//fmt.Printf("PollCount %s,\n", PollCount.Value)
		//fmt.Printf("RandomValue %s\n", RandomValue.Value)
	}
}
