package lbcfg

import (
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/digitalocean/godo"
)

// DropletStateInfo ...
type DropletStateInfo struct {
	Name             string
	IP               string
	State            bool
	FetchedTimestamp time.Time
}

// MapDropletStateInfo ...
type MapDropletStateInfo struct {
	Mutex              *sync.RWMutex
	Map                map[string]DropletStateInfo
	DigitaloceanClient *godo.Client
	DropletRegexp      *regexp.Regexp
	ChanChanges        chan bool
	NetworkType        string
}

// Init ...
func Init(
	apikey string,
	dropletRegexpString string,
	networkType string,
	nginxUpstreamConfigPath string,
	nginxBinaryPath string,
	apiPollDuration time.Duration,
	apiTimeout time.Duration,
	heartbeatPath string,
	heartbeatPollDuration time.Duration,
	heartbeatTimeout time.Duration,
	heartbeatThreshold int,
) {
	log.Printf("Starting...\n")

	digitaloceanClient := GetDigitaloceanClient(apikey)
	dropletRegexp, err := regexp.Compile(dropletRegexpString)
	if err != nil {
		panic("Error compiling DROPLET_REGEXP")
	}
	dropletStates := MapDropletStateInfo{
		Map:                make(map[string]DropletStateInfo),
		Mutex:              &sync.RWMutex{},
		DigitaloceanClient: digitaloceanClient,
		DropletRegexp:      dropletRegexp,
		ChanChanges:        make(chan bool, 1000),
		NetworkType:        networkType,
	}
	go dropletStates.DropletInfoLoop(apiPollDuration, apiTimeout)
	go dropletStates.DropletStateLoop(heartbeatPollDuration, heartbeatPath, heartbeatTimeout, heartbeatThreshold)
	go dropletStates.NginxConfiguratorLoop(nginxUpstreamConfigPath, nginxBinaryPath)
	go dropletStates.HeartbeatListener(heartbeatPath)

	select {}
}
