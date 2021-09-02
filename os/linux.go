package os

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/*
Public function returning score on VM is detected.
Scores Higher than 50 are highly likely to be a VM
*/
func VirtualMachineScore() int {
	score := 0

	// Check DMI Table for VM Entries
	if checkDMITable() {
		score + 50
		fmt.Println("DMI Table (/sys/class/dmi/id/*)")
	}

	// See if Kernel Detects a Hypervisor
	if checkKernelRingBuffer() {
		score + 50
		fmt.Println("Kernel Ring Buffer (/dev/kmsg)")
	}

	// Check for Hypervisor Flag or User Mode Linux
	if checkCPUInfo() {
		score + 50
		fmt.Println("CPU Vendor (/proc/cpuinfo)")
	}

	// Look for VM Tools in Modules
	if checkKernelModules() {
		score + 20
		fmt.Println("Kernel module (/proc/modules)")
	}

	// Check for Xen Tools
	if checkXenProcFile() {
		score + 20
		fmt.Println("Xen proc file (/proc/xen)")
	}

	// Some Distros Contain This File
	if checkHypervisorType() {
		score + 10
		fmt.Println("Hypervisor type file (/sys/hypervisor/type)")
	}

	// Some Distros Contain VM Marker in proc
	if checkSysInfo() {
		score + 10
		fmt.Println("System Information (/proc/sysinfo)")
	}

	// Fairly Rare Check
	if checkDeviceTree() {
		score + 10
		fmt.Println("VM device tree (/proc/device-tree)")
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
	if DoesFileContain(file, "vboxguest") {
		return true
	}
	if DoesFileContain(file, "xenfs") {
		return true
	}
	if DoesFileContain(file, "qemu") {
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
		PrintError(err)
		return false
	}

	for _, dmiEntry := range dmiFiles {
		if !dmiEntry.Mode().IsRegular() {
			continue
		}

		dmiContent, err := ioutil.ReadFile(dmiPath + dmiEntry.Name())

		if err != nil {
			PrintError(err)
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
		PrintError(err)
		return false
	}

	defer file.Close()

	// Set a read timeout because otherwise reading kmsg (which is a character device) will block
	if err = file.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		PrintError(err)
		return false
	}

	return DoesFileContain(file, "Hypervisor detected")
}

// Look in CPUInfo for the hypervisor flag
// Checks if UML is being used - https://en.wikipedia.org/wiki/User-mode_Linux
func checkCPUInfo() bool {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		PrintError(err)
		return false
	}
	defer file.Close()

	if DoesFileContain(file, "hypervisor") {
		return true
	}
	if DoesFileContain(file, "User Mode Linux") {
		return true
	}

	return false
}

// Some GNU/Linux distributions expose /proc/sysinfo containing potential VM info
// https://www.ibm.com/support/knowledgecenter/en/linuxonibm/com.ibm.linux.z.lhdd/lhdd_t_sysinfo.html
func checkSysInfo() bool {
	file, err := os.Open("/proc/sysinfo")
	if err != nil {
		PrintError(err)
		return false
	}
	defer file.Close()

	return DoesFileContain(file, "VM00")
}

// Some virtualization technologies can be detected using /proc/device-tree
func checkDeviceTree() bool {
	deviceTreeBase := "/proc/device-tree"

	if DoesFileExist(deviceTreeBase + "/hypervisor/compatible") {
		return true
	}
	if DoesFileExist(deviceTreeBase + "/fw-cfg") {
		return true
	}

	return false
}

// Some virtualization technologies can be detected using /proc/type
func checkHypervisorType() bool {
	return DoesFileExist("/sys/hypervisor/type")
}

// Xen can be detected thanks to /proc/xen
func checkXenProcFile() bool {
	return DoesFileExist("/proc/xen")
}
