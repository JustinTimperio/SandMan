package detection

import (
	"fmt"
	"net"
	"strings"
)

// SandboxMac is used to check if the environment's MAC address
// matches standard MAC adddresses of virtualized environments.
func SandboxMac() bool {
	sandbox_macs := []string{
		`00:0C:29`, // VMWare
		`00:1C:14`, // VMWare
		`00:50:56`, // VMWare
		`00:05:69`, // VMWare
		`08:00:27`, // VirtualBox
		`00:0F:4F`, // Oracle VirtualBox
		`02:42:ac`, // Docker
		"00:1C:42", // Parallels
		"00:16:E3", // Citrix Xen (XCP-NG uses random UUID's)
		"00:15:5d", // Hyper-V

	}
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		for _, mac := range sandbox_macs {
			if strings.Contains(strings.ToLower(iface.HardwareAddr.String()), strings.ToLower(mac)) {
				fmt.Println(iface.HardwareAddr.String())
				return true
			}
		}
	}

	return false
}
