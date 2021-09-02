package detection

import (
	"runtime"

	"github.com/klauspost/cpuid"
)

// SandboxRam is used to check if the environment's
// RAM is less than a given size.
func SandboxRam() bool {
	ram_mb := 1024
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	rmb := uint64(ram_mb)
	ram := m.TotalAlloc / 1024 / 1024

	return ram <= rmb
}

func VMCPUTest() bool {
	if cpuid.CPU.VM() {
		return true
	}
	return false
}
