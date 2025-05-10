//go:build windows

package pid

import (
	"fmt"
	"runtime"
	"time"

	"golang.org/x/sys/windows"

	"github.com/ravoni4devs/syspector/internal/common"
)

func getStat(pidNumber int, interval time.Duration) (PidStat, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pidNumber))
	if err != nil {
		return PidStat{}, err
	}
	defer windows.CloseHandle(handle)

	var creationTime, exitTime, kernelTime1, userTime1 windows.Filetime
	err = windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime1, &userTime1)
	if err != nil {
		return PidStat{}, err
	}

	time.Sleep(interval)

	var kernelTime2, userTime2 windows.Filetime
	_ = windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime2, &userTime2)

	// Subtração direta entre os durations
	deltaUser := filetimeToDuration(userTime2) - filetimeToDuration(userTime1)
	deltaKernel := filetimeToDuration(kernelTime2) - filetimeToDuration(kernelTime1)

	totalTime := deltaUser + deltaKernel
	usage := (totalTime.Seconds() / interval.Seconds()) * 100.0 / float64(runtime.NumCPU())

	return PidStat{PID: pidNumber, CpuPercent: common.ParseFloat(fmt.Sprintf("%.2f", usage))}, nil
}

func filetimeToDuration(ft windows.Filetime) time.Duration {
	return time.Duration(ft.HighDateTime)<<32 + time.Duration(ft.LowDateTime)
}
