package system

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ravoni4devs/syspector/internal/common"
)

func Uptime() (float64, error) {
	// Darwin (macOS) pode não ter syscall.Sysinfo implementado
	// Alternativa usando sysctl
	type timeval struct {
		Sec  int32
		Usec int32
	}

	invoke := common.Invoke{}
	b, err := invoke.Command("sysctl", "-n", "kern.boottime")
	if err != nil {
		return 0, err
	}
	parts := strings.Split(string(b), "=")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid format of kern.bootime: %s", string(b))
	}
	secPart := strings.Split(parts[1], ",")[0]
	secStr := strings.TrimSpace(secPart)
	sec, err := strconv.ParseInt(secStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert to secs: %s", secStr)
	}
	bootTime := time.Unix(sec, 0)
	uptime := time.Since(bootTime).Seconds()
	return uptime, nil
}
