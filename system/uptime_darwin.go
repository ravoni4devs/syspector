package system

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

func Uptime() (float64, error) {
	// Darwin (macOS) pode não ter syscall.Sysinfo implementado
	// Alternativa usando sysctl
	type timeval struct {
		Sec  int32
		Usec int32
	}

	var bootTime timeval
	mib := []int32{1 /* CTL_KERN */, 21 /* KERN_BOOTTIME */}
	size := uintptr(unsafe.Sizeof(bootTime))

	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		2,
		uintptr(unsafe.Pointer(&bootTime)),
		uintptr(unsafe.Pointer(&size)),
		0,
		0,
	)

	if errno != 0 {
		return 0, fmt.Errorf("falha ao obter boot time: %v", errno)
	}

	uptime := time.Now().Sub(time.Unix(int64(bootTime.Sec), int64(bootTime.Usec)*1000))
	return uptime.Seconds(), nil
}
