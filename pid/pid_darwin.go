//go:build darwin

package pid

import (
	"fmt"
	"runtime"
	"syscall"
	"time"

	"github.com/ravoni4devs/syspector/internal/common"
)

func getStat(pidNumber int, interval time.Duration) (PidStat, error) {
    var r1, r2 syscall.Rusage
    err := syscall.Getrusage(syscall.RUSAGE_SELF, &r1)
    if err != nil {
        return PidStat{}, err
    }

    time.Sleep(interval)

    err = syscall.Getrusage(syscall.RUSAGE_SELF, &r2)
    if err != nil {
        return PidStat{}, err
    }

    user1 := time.Duration(r1.Utime.Sec)*time.Second + time.Duration(r1.Utime.Usec)*time.Microsecond
    sys1 := time.Duration(r1.Stime.Sec)*time.Second + time.Duration(r1.Stime.Usec)*time.Microsecond
    user2 := time.Duration(r2.Utime.Sec)*time.Second + time.Duration(r2.Utime.Usec)*time.Microsecond
    sys2 := time.Duration(r2.Stime.Sec)*time.Second + time.Duration(r2.Stime.Usec)*time.Microsecond

    deltaTime := user2 - user1 + sys2 - sys1

    usage := (deltaTime.Seconds() / interval.Seconds()) * 100.0 / float64(runtime.NumCPU())

    return PidStat{PID: syscall.Getpid(), CpuPercent: common.ParseFloat(fmt.Sprintf("%.2f", usage))}, nil
}

