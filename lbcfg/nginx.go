package lbcfg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// NginxConfiguratorLoop ...
func (dropletStates MapDropletStateInfo) NginxConfiguratorLoop(nginxUpstreamConfigPath string, nginxBinaryPath string) {
	for {
		select {
		case <-dropletStates.ChanChanges:
			dropletStates.Mutex.RLock()
			var buffer string
			buffer += fmt.Sprintf("upstream backend {\n")
			keepalive := len(dropletStates.Map) * 10
			for _, state := range dropletStates.Map {
				buffer += fmt.Sprintf("\t# %s\n", state.Name)
				if state.State {
					buffer += fmt.Sprintf("\tserver %s:80;\n\t\n", state.IP)
				} else {
					buffer += fmt.Sprintf("\tserver %s:80 down;\n\t\n", state.IP)
				}
			}
			buffer += fmt.Sprintf("\tleast_conn;\n\tkeepalive %d;\n}\n", keepalive)
			dropletStates.Mutex.RUnlock()

			configFile, err := os.OpenFile(nginxUpstreamConfigPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("Error opening nginx upstream config for write: %v\n", err)
				continue
			}
			configFile.Write([]byte(buffer))
			configFile.Close()

			cmd := exec.Command(nginxBinaryPath, "-s", "reload")
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error reloading nginx: %v\n", err)
				continue
			}
			fmt.Printf("NGINX: %s\n", string(out))
		}
	}
}
