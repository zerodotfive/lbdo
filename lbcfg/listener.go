package lbcfg

import (
	"fmt"
	"log"
	"net/http"
)

// HeartbeatListener ...
func (dropletStates MapDropletStateInfo) HeartbeatListener(heartbeatPath string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		dropletStates.Mutex.RLock()
		defer func() {
			dropletStates.Mutex.RUnlock()
		}()
		for _, droplet := range dropletStates.Map {
			if droplet.State {
				fmt.Fprintf(w, "true")
				return
			}
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "false")
	}
	http.HandleFunc(heartbeatPath, handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
