package SandMan

import "SandMan/detection"

func Evade() {

}

func Detect() int {
	score := detection.VirtualMachineScore()
	return score
}
