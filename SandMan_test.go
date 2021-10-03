package SandMan_test

import (
	"SandMan"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {

	//fmt.Println("Hardware VM score:", SandMan.DetectHardware())
	fmt.Println("VM Artifact score:", SandMan.DetectVMArtifacts())
}
