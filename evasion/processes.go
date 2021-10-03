package evasion

import (
	"SandMan/shared"

	ps "github.com/mitchellh/go-ps"
)

// SandboxProc checks if there are processes that indicate
// a virtualized environment.
func SandboxProc() bool {
	sandbox_processes := []string{
		`vmsrvc`,
		`tcpview`,
		`wireshark`,
		`visual basic`,
		`fiddler`,
		`vmware`,
		`vbox`,
		`process explorer`,
		`autoit`,
		`vboxtray`,
		`vmtools`,
		`vmrawdsk`,
		`vmusbmouse`,
		`vmvss`,
		`vmscsi`,
		`vmxnet`,
		`vmx_svga`,
		`vmmemctl`,
		`df5serv`,
		`vboxservice`,
		`vmhgfs`,
	}
	p, _ := processes()
	for _, name := range p {
		if shared.ContainsAny(name, sandbox_processes) {
			return true
		}
	}
	return false
}

// processes returns a map of a PID to its respective process name.
func processes() (map[int]string, error) {
	prs := make(map[int]string)
	processList, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	for x := range processList {
		process := processList[x]
		prs[process.Pid()] = process.Executable()
	}

	return prs, nil
}
