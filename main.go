package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/zerodotfive/lbdo/lbcfg"
)

const (
	defaultNetworkType             = "private"
	defaultNginxUpstreamConfigPath = "/tmp/upstream.conf"
	defaultNginxBinaryPath         = "/usr/sbin/nginx"
	defaultAPIPollDuration         = time.Second * 30
	defaultAPITimeout              = time.Second * 10
	defaultHeartbeatPath           = "/"
	defaultHeartbeatPollDuration   = time.Second * 10
	defaultHeartbeatTimeout        = time.Second * 2
	defaultHeartbeatThreshold      = 3
)

func help(
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
	fmt.Printf("ENV list:\n")
	fmt.Printf("APIKEY\n\tDigitalOcean api key\n\n")
	fmt.Printf("DROPLET_REGEXP\n\tDroplet regexp \"%s\"\n\n", dropletRegexpString)
	fmt.Printf("NETWORK_TYPE\n\tNetwork type \"%s\". Can be \"private\" or \"public\". By default: \"%s\"\n\n", networkType, defaultNetworkType)
	fmt.Printf("NGINX_UPSTREAM_CONFIG_PATH\n\tNginx upstream config path \"%s\". By default: \"%s\"\n\n", nginxUpstreamConfigPath, defaultNginxUpstreamConfigPath)
	fmt.Printf("NGINX_BINARY_PATH\n\tNginx binary path \"%s\". By default: \"%s\"\n\n", nginxBinaryPath, defaultNginxBinaryPath)
	fmt.Printf("API_POLL_DURATION\n\tApi poll duration \"%s\". By default: \"%s\"\n\n", apiPollDuration, defaultAPIPollDuration)
	fmt.Printf("API_TIMEOUT\n\tApi timeout \"%s\". By default: \"%s\"\n\n", apiTimeout, defaultAPITimeout)
	fmt.Printf("HEARTBEAT_PATH\n\tHeartbeat path \"%s\". By default: \"%s\"\n\n", heartbeatPath, defaultHeartbeatPath)
	fmt.Printf("HEARTBEAT_POLL_DURATION\n\tHeartbeat poll duration \"%s\". By default: \"%s\"\n\n", heartbeatPollDuration, defaultHeartbeatPollDuration)
	fmt.Printf("HEARTBEAT_TIMEOUT\n\tHeartbeat timeout \"%s\". By default: \"%s\"\n\n", heartbeatTimeout, defaultHeartbeatTimeout)
	fmt.Printf("HEARTBEAT_THRESHOLD\n\tHeartbeat threshold \"%d\". By default: \"%d\"\n\n", heartbeatThreshold, defaultHeartbeatThreshold)
}

func main() {
	viper.AutomaticEnv()

	apikey := viper.GetString("APIKEY")

	dropletRegexpString := viper.GetString("DROPLET_REGEXP")

	networkType := viper.GetString("NETWORK_TYPE")
	if len(networkType) == 0 {
		networkType = defaultNetworkType
	}

	nginxConfigPath := viper.GetString("NGINX_UPSTREAM_CONFIG_PATH")
	if len(nginxConfigPath) == 0 {
		nginxConfigPath = defaultNginxUpstreamConfigPath
	}

	nginxBinaryPath := viper.GetString("NGINX_BINARY_PATH")
	if len(nginxBinaryPath) == 0 {
		nginxBinaryPath = defaultNginxBinaryPath
	}

	apiPollDuration := viper.GetDuration("API_POLL_DURATION")
	if apiPollDuration == 0 {
		apiPollDuration = defaultAPIPollDuration
	}

	apiTimeout := viper.GetDuration("API_TIMEOUT")
	if apiTimeout == 0 {
		apiTimeout = defaultAPITimeout
	}

	heartbeatPath := viper.GetString("HEARTBEAT_PATH")
	if len(heartbeatPath) == 0 {
		heartbeatPath = defaultHeartbeatPath
	}

	heartbeatPollDuration := viper.GetDuration("HEARTBEAT_POLL_DURATION")
	if heartbeatPollDuration == 0 {
		heartbeatPollDuration = defaultHeartbeatPollDuration
	}

	heartbeatTimeout := viper.GetDuration("HEARTBEAT_TIMEOUT")
	if heartbeatTimeout == 0 {
		heartbeatTimeout = defaultHeartbeatTimeout
	}

	heartbeatThreshold := viper.GetInt("HEARTBEAT_THRESHOLD")
	if heartbeatThreshold == 0 {
		heartbeatThreshold = defaultHeartbeatThreshold
	}

	if len(dropletRegexpString) == 0 {
		help(
			apikey,
			dropletRegexpString,
			networkType,
			nginxConfigPath,
			nginxBinaryPath,
			apiPollDuration,
			apiTimeout,
			heartbeatPath,
			heartbeatPollDuration,
			heartbeatTimeout,
			heartbeatThreshold,
		)
		panic("Please, set DROPLET_REGEXP env\n")
	}
	if len(apikey) == 0 {
		help(
			apikey,
			dropletRegexpString,
			networkType,
			nginxConfigPath,
			nginxBinaryPath,
			apiPollDuration,
			apiTimeout,
			heartbeatPath,
			heartbeatPollDuration,
			heartbeatTimeout,
			heartbeatThreshold,
		)
		panic("Please, set APIKEY env\n")
	}
	if networkType != "private" && networkType != "public" {
		help(
			apikey,
			dropletRegexpString,
			networkType,
			nginxConfigPath,
			nginxBinaryPath,
			apiPollDuration,
			apiTimeout,
			heartbeatPath,
			heartbeatPollDuration,
			heartbeatTimeout,
			heartbeatThreshold,
		)
		panic("Please, set NETWORK_TYPE to \"private\" or  \"public\"\n")
	}

	if (heartbeatTimeout.Seconds() * float64(heartbeatThreshold)) >= heartbeatPollDuration.Seconds() {
		panic("Error: HEARTBEAT_TIMEOUT * HEARTBEAT_THRESHOLD should be less than HEARTBEAT_POLL_DURATION")
	}

	help(
		apikey,
		dropletRegexpString,
		networkType,
		nginxConfigPath,
		nginxBinaryPath,
		apiPollDuration,
		apiTimeout,
		heartbeatPath,
		heartbeatPollDuration,
		heartbeatTimeout,
		heartbeatThreshold,
	)

	lbcfg.Init(
		apikey,
		dropletRegexpString,
		networkType,
		nginxConfigPath,
		nginxBinaryPath,
		apiPollDuration,
		apiTimeout,
		heartbeatPath,
		heartbeatPollDuration,
		heartbeatTimeout,
		heartbeatThreshold,
	)
}
