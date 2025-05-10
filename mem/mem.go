package mem

import (
	"encoding/json"

	"github.com/ravoni4devs/syspector/internal/common"
)

var invoke common.Invoker = common.Invoke{}

// Memory usage statistics. Total, Available and Used contain numbers of bytes
// for human consumption.
//
// The other fields in this struct contain kernel specific values.
type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`

	// OS X / BSD specific numbers:
	// http://www.macyourself.com/2010/02/17/what-is-free-wired-active-and-inactive-system-memory-ram/
	Active   uint64 `json:"active,omitempty"`
	Inactive uint64 `json:"inactive,omitempty"`
	Wired    uint64 `json:"wired,omitempty"`

	// FreeBSD specific numbers:
	// https://reviews.freebsd.org/D8467
	Laundry uint64 `json:"laundry,omitempty"`

	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	// https://www.kernel.org/doc/Documentation/vm/transhuge.txt
	Buffers        uint64 `json:"buffers,omitempty"`
	Cached         uint64 `json:"cached,omitempty"`
	WriteBack      uint64 `json:"writeBack,omitempty"`
	Dirty          uint64 `json:"dirty,omitempty"`
	WriteBackTmp   uint64 `json:"writeBackTmp,omitempty"`
	Shared         uint64 `json:"shared,omitempty"`
	Slab           uint64 `json:"slab,omitempty"`
	Sreclaimable   uint64 `json:"sreclaimable,omitempty"`
	Sunreclaim     uint64 `json:"sunreclaim,omitempty"`
	PageTables     uint64 `json:"pageTables,omitempty"`
	SwapCached     uint64 `json:"swapCached,omitempty"`
	CommitLimit    uint64 `json:"commitLimit,omitempty"`
	CommittedAS    uint64 `json:"committedAS,omitempty"`
	HighTotal      uint64 `json:"highTotal,omitempty"`
	HighFree       uint64 `json:"highFree,omitempty"`
	LowTotal       uint64 `json:"lowTotal,omitempty"`
	LowFree        uint64 `json:"lowFree,omitempty"`
	SwapTotal      uint64 `json:"swapTotal,omitempty"`
	SwapFree       uint64 `json:"swapFree,omitempty"`
	Mapped         uint64 `json:"mapped,omitempty"`
	VmallocTotal   uint64 `json:"vmallocTotal,omitempty"`
	VmallocUsed    uint64 `json:"vmallocUsed,omitempty"`
	VmallocChunk   uint64 `json:"vmallocChunk,omitempty"`
	HugePagesTotal uint64 `json:"hugePagesTotal,omitempty"`
	HugePagesFree  uint64 `json:"hugePagesFree,omitempty"`
	HugePagesRsvd  uint64 `json:"hugePagesRsvd,omitempty"`
	HugePagesSurp  uint64 `json:"hugePagesSurp,omitempty"`
	HugePageSize   uint64 `json:"hugePageSize,omitempty"`
	AnonHugePages  uint64 `json:"anonHugePages,omitempty"`
}

type SwapMemoryStat struct {
	Total       uint64  `json:"total,omitempty"`
	Used        uint64  `json:"used,omitempty"`
	Free        uint64  `json:"free,omitempty"`
	UsedPercent float64 `json:"usedPercent,omitempty"`
	Sin         uint64  `json:"sin,omitempty"`
	Sout        uint64  `json:"sout,omitempty"`
	PgIn        uint64  `json:"pgIn,omitempty"`
	PgOut       uint64  `json:"pgOut,omitempty"`
	PgFault     uint64  `json:"pgFault,omitempty"`

	// Linux specific numbers
	// https://www.kernel.org/doc/Documentation/cgroup-v2.txt
	PgMajFault uint64 `json:"pgMajFault,omitempty"`
}

func (m VirtualMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func (m SwapMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

type SwapDevice struct {
	Name      string `json:"name"`
	UsedBytes uint64 `json:"usedBytes"`
	FreeBytes uint64 `json:"freeBytes"`
}

func (m SwapDevice) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}
