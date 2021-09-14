package detection

import (
	"SandMan/shared"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/*
VirtualMachineScore is a public function returning a VM likelyhood rating.
Different methods are given different weights but typically scores >= 50 are highly likely to be a VM
*/
func VirtualMachineScore() int {
	score := 0

	// Check DMI Table for VM Entries
	if checkDMITable() {
		score = score + 50
		fmt.Println("VM detected in DMI Table (/sys/class/dmi/id/*)")
	}

	// See if Kernel Detects a Hypervisor
	if checkKernelRingBuffer() {
		score = score + 50
		fmt.Println("VM detected in Kernel Ring Buffer (/dev/kmsg)")
	}

	// Check for Hypervisor Flag or User Mode Linux
	if checkCPUInfo() {
		score = score + 50
		fmt.Println("VM detected in CPU Vendor (/proc/cpuinfo)")
	}

	// Check the Device Tree for VM markers
	if checkDeviceTree() {
		score = score + 50
		fmt.Println("VM detected in device tree (/proc/*)")
	}

	// Look for VM Tools in Modules
	if checkKernelModules() {
		score = score + 50
		fmt.Println("VM detected in kernel module (/proc/modules)")
	}

	// Some Distros Contain VM Marker in proc
	if checkSysInfo() {
		score = score + 10
		fmt.Println("VM detected in System Information (/proc/sysinfo)")
	}

	return score
}

// Detect Virtual Machine Modules in proc
func checkKernelModules() bool {

	file, err := os.Open("/proc/modules")
	if err != nil {
		fmt.Println(err)
		return false
	}

	// TODO: Add More Entries
	if shared.DoesFileContain(file, "vboxguest") {
		return true
	}
	if shared.DoesFileContain(file, "xenfs") {
		return true
	}
	if shared.DoesFileContain(file, "qemu") {
		return true
	}
	if shared.DoesFileContain(file, "vmw_vsock") {
		return true
	}

	return false
}

//Checks if the DMI table contains vendor strings of known VMs.
func checkDMITable() bool {

	blacklistDMI := []string{
		// Must be lowercase
		"innotek",
		"virtualbox",
		"vbox",
		"kvm",
		"qemu",
		"vmware",
		"vmw",
		"oracle",
		"xen",
		"bochs",
		"parallels",
		"bhyve",
	}

	dmiPath := "/sys/class/dmi/id/"
	dmiFiles, err := ioutil.ReadDir(dmiPath)

	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, dmiEntry := range dmiFiles {
		if !dmiEntry.Mode().IsRegular() {
			continue
		}

		dmiContent, err := ioutil.ReadFile(dmiPath + dmiEntry.Name())

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, entry := range blacklistDMI {
			// Lowercase comparison to prevent false negatives
			if bytes.Contains(bytes.ToLower(dmiContent), []byte(entry)) {
				return true
			}
		}
	}

	return false
}

// Checks printk messages to see if Linux detected an hypervisor.
// https://github.com/torvalds/linux/blob/31cc088a4f5d83481c6f5041bd6eb06115b974af/arch/x86/kernel/cpu/hypervisor.c#L79
func checkKernelRingBuffer() bool {

	file, err := os.Open("/dev/kmsg")

	if err != nil {
		fmt.Println(err)
		return false
	}

	defer file.Close()

	// Set a read timeout because otherwise reading kmsg (which is a character device) will block
	if err = file.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		fmt.Println(err)
		return false
	}

	return shared.DoesFileContain(file, "Hypervisor detected")
}

// Look in CPUInfo for the hypervisor flag
// Checks if UML is being used - https://en.wikipedia.org/wiki/User-mode_Linux
func checkCPUInfo() bool {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	if shared.DoesFileContain(file, "hypervisor") {
		return true
	}
	if shared.DoesFileContain(file, "User Mode Linux") {
		return true
	}

	return false
}

// Run a bunch of file system checks for various VM markers in /proc/
func checkDeviceTree() bool {
	if shared.DoesPathExist("/proc/sys/xen") {
		return true
	}
	if shared.DoesPathExist("/proc/xen") {
		return true
	}
	if shared.DoesPathExist("/proc/device-tree/hypervisor/compatible") {
		return true
	}
	if shared.DoesPathExist("/proc/device-tree/fw-cfg") {
		return true
	}
	if shared.DoesPathExist("/sys/hypervisor/type") {
		return true
	}

	return false
}

// Some GNU/Linux distributions expose /proc/sysinfo containing potential VM info
// https://www.ibm.com/support/knowledgecenter/en/linuxonibm/com.ibm.linux.z.lhdd/lhdd_t_sysinfo.html
func checkSysInfo() bool {

	if shared.DoesPathExist("/proc/sysinfo") {
		file, err := os.Open("/proc/sysinfo")
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer file.Close()

		return shared.DoesFileContain(file, "VM00")
	}
	return false
}
