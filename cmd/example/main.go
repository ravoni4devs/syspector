package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ravoni4devs/syspector"
	"github.com/ravoni4devs/syspector/cpu"
	"github.com/ravoni4devs/syspector/docker"
	"github.com/ravoni4devs/syspector/mem"
)

func main() {
	var (
		httpPort  string
		usePid    bool
		useDocker bool
		useHttp   bool
		useMemory bool
		useCpu    bool
	)
	flag.BoolVar(&usePid, "pid", false, "Print stats for current PID")
	flag.StringVar(&httpPort, "port", "8080", "Http port used together with -http param")
	flag.BoolVar(&useDocker, "docker", false, "Print docker container stats")
	flag.BoolVar(&useHttp, "http", false, "Expose API rest with stats")
	flag.BoolVar(&useMemory, "memory", false, "Print memory stats")
	flag.BoolVar(&useCpu, "cpu", false, "Print cpu stats")

	flag.Parse()

	if useDocker {
		printDockerStats()
		return
	}

	if useHttp {
		runHttpServer(httpPort)
		return
	}

	if useMemory {
		printMemoryStats()
		return
	}

	if useCpu {
		printCpuStats()
		return
	}

	if usePid {
		printPidStats(os.Getpid())
		return
	}

	stats, err := syspector.New().Stats()
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}
	fmt.Println(prettyJSON(stats))
}

func printDockerStats() {
	v, _ := docker.VirtualMemory()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	fmt.Println(v)
}

func runHttpServer(port string) {
	port = ":" + port
	collector := syspector.New()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stats, _ := collector.Stats()
		data := fmt.Sprintf(`{"data": "%s"}`, prettyJSON(stats))
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"message": "Method not allowed"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})
	log.Println("Listening", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func printMemoryStats() {
	v, _ := mem.GetStat()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	fmt.Println(v)
}

func printCpuStats() {
	usage, _ := cpu.Percent(time.Second*1, false)
	fmt.Println("CPU usage:", usage)
}

func printPidStats(pid int) {
	fmt.Println("PID=", pid)
	collector := syspector.New()
	interval := 1 * time.Second
	go func() {
		for {
			stats, _ := collector.GetStatsByPID(pid)
			fmt.Println(prettyJSON(stats))
			<-time.After(interval)
			runtime.GC()
		}
	}()
	_ = make([]byte, 1*(1024*1024))
	<-time.After(time.Duration(math.MaxInt64))
}

func prettyJSON(i any) string {
	if i == nil {
		return ""
	}
	b, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}
