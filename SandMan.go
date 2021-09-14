package SandMan

import (
	"SandMan/detection"
	"fmt"
)

func Evade() {

}

func DetectHardware() int {
	score := 0

	// Start by checking time compression
	tc := detection.SandboxTimeCompression()
	if tc {
		fmt.Println("The OS or hardware is fucking with time!")
		score++
	}

	mac := detection.SandboxMac()
	if mac {
		fmt.Println("Sandbox Mac detected!")
		score++
	}

	smolRam := detection.SandboxRam()
	if smolRam {
		fmt.Println("Small amount of Ram detected!")
		score++
	}

	vmCpu := detection.VMCPUTest()
	if vmCpu {
		fmt.Println("VM CPU detected!")
		score++
	}

	// Double check time compression just incase sys calls trigger speed up
	tc = detection.SandboxTimeCompression()
	if tc {
		fmt.Println("The OS or hardware is fucking with time!")
		score++
	}

	return score
}

func DetectVMArtifacts() int {
	score := detection.VirtualMachineScore()
	return score
}
