# Syspector

Syspector is a cross-platform library designed to gather system metrics, including CPU and memory usage, system information, and
application resource consumption. It can also retrieve stats from Docker containers. The library supports Linux,
Windows, and macOS, and can be used as a command-line tool or integrated into your own applications.

## Features

- **Memory Usage:** Get detailed memory stats such as total, free, and used memory, along with the percentage of memory in use.
- **CPU Usage:** Retrieve the current CPU usage for your application or system.
- **System Information:** Gather various system metrics, such as operating system details and architecture.
- **Application Consumption:** Fetch resource usage of the current application, including CPU and memory.
- **Docker Containers:** Fetch stats for Docker containers, including memory and CPU usage.
- **Cross-Platform:** Works on Linux, Windows, and macOS (Darwin).

## Installation

```sh
go get github.com/ravoni4devs/syspector
```

## Usage

The following is an example of how you can use the syspector library to get system and application stats.

```go
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
```

**Available Flags**

- **pid:** Get stats for the current process ID (PID).
- **port:** Set the HTTP port for the REST API (default:** 8080).
- **docker:** Fetch stats for Docker containers.
- **http:** Expose the stats via a REST API.
- **memory:** Print memory stats (total, free, used percentage).
- **cpu:** Print CPU stats (percentage of CPU usage).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
