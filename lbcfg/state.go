package lbcfg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// getDropletState ...
func (dropletStates MapDropletStateInfo) getDropletState(name string, heartbeatPath string, heartbeatTimeout time.Duration, heartbeatThreshold int) {
	dropletStates.Mutex.RLock()
	ip := dropletStates.Map[name].IP
	oldState := dropletStates.Map[name].State
	fetchedTimestamp := dropletStates.Map[name].FetchedTimestamp
	dropletStates.Mutex.RUnlock()
	newState := oldState
	failCount := 0
	successCount := 0
	defer func() {
		if failCount == heartbeatThreshold {
			newState = false
		}
		if successCount == heartbeatThreshold {
			newState = true
		}
		if oldState != newState {
			dropletStates.Mutex.Lock()
			dropletStates.Map[name] = DropletStateInfo{
				Name:             name,
				IP:               ip,
				State:            newState,
				FetchedTimestamp: fetchedTimestamp,
			}
			dropletStates.ChanChanges <- true
			dropletStates.Mutex.Unlock()
		}
	}()
	for c := 1; c <= heartbeatThreshold; c++ {
		checkCtx, cancel := context.WithTimeout(context.Background(), heartbeatTimeout)
		request, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", dropletStates.Map[name].IP, heartbeatPath), nil)
		if err != nil {
			failCount++
			if oldState || successCount > 0 {
				log.Printf("%s check fail #%d\n", name, failCount)
			}
			cancel()
			continue
		}
		request = request.WithContext(checkCtx)
		client := http.DefaultClient
		response, err := client.Do(request)
		if err != nil {
			failCount++
			if oldState || successCount > 0 {
				log.Printf("%s check fail #%d\n", name, failCount)
			}
			cancel()
			continue
		}
		if response.StatusCode == 200 {
			successCount++
			if !oldState || failCount > 0 {
				log.Printf("%s check success #%d\n", name, successCount)
			}
			cancel()
			continue
		}
		failCount++
		if oldState || successCount > 0 {
			log.Printf("%s check fail #%d\n", name, failCount)
		}
		cancel()
	}
	return
}

// DropletStateLoop ...
func (dropletStates MapDropletStateInfo) DropletStateLoop(
	heartbeatPollDuration time.Duration,
	heartbeatPath string,
	heartbeatTimeout time.Duration,
	heartbeatThreshold int,
) {
	for {
		dropletStates.Mutex.Lock()
		for name, state := range dropletStates.Map {
			// Remove deprecated
			if time.Now().Sub(state.FetchedTimestamp) >= time.Hour {
				delete(dropletStates.Map, name)
				dropletStates.ChanChanges <- true
				continue
			}
			go dropletStates.getDropletState(name, heartbeatPath, heartbeatTimeout, heartbeatThreshold)
		}
		dropletStates.Mutex.Unlock()
		time.Sleep(heartbeatPollDuration)
	}
}
