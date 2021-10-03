package detection

import (
	"github.com/klauspost/cpuid"
	"github.com/pbnjay/memory"
)

// SandboxRam is used to check if the environment's
// RAM is less than a given size.
func SandboxRam() bool {
	ramSmol := 1073741824
	ram := int(memory.TotalMemory())
	return ram <= ramSmol
}

func VMCPUTest() bool {
	return cpuid.CPU.VM()
}
