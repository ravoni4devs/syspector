package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ravoni4devs/syspector/internal/common"
)

const (
	procBasePath    = "/proc"
	uptimeFilePath  = "uptime"
	meminfoFilePath = "meminfo"
	cgroupFilePath  = "1/cgroup"
)

type SystemStat struct {
	Uptime       float64 `json:"uptime"`
	CPUs         int     `json:"cpus"`
	OSFamily     string  `json:"os"`
	Architecture string  `json:"arch"`
	Version      string  `json:"version"`
	Distro       string  `json:"distro"`
	Container    string  `json:"container"`
	Virtualized  bool    `json:"virtualized"`
}

func GetStat() (SystemStat, error) {
	var stat = SystemStat{
		CPUs: NumCPU(),
	}
	uptime, err := Uptime()
	if err != nil {
		return stat, err
	}
	stat.Uptime = uptime
	stat.OSFamily = runtime.GOOS
	stat.Architecture = runtime.GOARCH
	stat.Version = getOSVersion()
	if stat.OSFamily == "linux" {
		stat.Distro = getLinuxDistro()
		stat.Container = detectContainer()
		stat.Virtualized = isVirtualized()
	}
	return stat, nil
}

func GetOSInfo() (SystemStat, error) {
	var stat = SystemStat{}
	stat.OSFamily = runtime.GOOS
	stat.Architecture = runtime.GOARCH
	stat.Version = getOSVersion()
	if stat.OSFamily == "linux" {
		stat.Distro = getLinuxDistro()
		stat.Container = detectContainer()
		stat.Virtualized = isVirtualized()
	}
	return stat, nil
}

func Uptime() (float64, error) {
	filename := fmt.Sprintf("%s/%s", procBasePath, uptimeFilePath)
	data, err := common.ReadFileNoStat(filename)
	if err != nil {
		return 0, err
	}
	return common.ParseFloat(strings.Split(string(data), " ")[0]), nil
}

func NumCPU() int {
	return runtime.NumCPU()
}
