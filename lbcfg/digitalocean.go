package lbcfg

import (
	"context"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// GetDigitaloceanClient ...
func GetDigitaloceanClient(apikey string) *godo.Client {
	authCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token := &tokenSource{
		AccessToken: apikey,
	}
	oauthClient := oauth2.NewClient(authCtx, token)
	return godo.NewClient(oauthClient)
}

// DropletInfoLoop ...
func (dropletStates MapDropletStateInfo) DropletInfoLoop(apiPollDuration time.Duration, apiTimeout time.Duration) error {
MAIN_LOOP:
	for {
		opt := &godo.ListOptions{}
	PAGE_LOOP:
		for {
			dropletCtx, cancel := context.WithTimeout(context.Background(), apiTimeout)
			droplets, resp, err := dropletStates.DigitaloceanClient.Droplets.List(dropletCtx, opt)
			cancel()
			if err != nil {
				log.Printf("%v\n", err)
				time.Sleep(apiPollDuration)
				continue MAIN_LOOP
			}
			for _, d := range droplets {
				matched := dropletStates.DropletRegexp.Match([]byte(d.Name))
				if matched {
					var ip string
					for _, n := range d.Networks.V4 {
						if n.Type == dropletStates.NetworkType {
							ip = n.IPAddress
						}
					}
					dropletStates.Mutex.Lock()
					if val, ok := dropletStates.Map[d.Name]; ok {
						dropletStates.Map[d.Name] = DropletStateInfo{
							IP:               val.IP,
							Name:             d.Name,
							State:            val.State,
							FetchedTimestamp: time.Now(),
						}
					} else {
						dropletStates.Map[d.Name] = DropletStateInfo{
							IP:               ip,
							Name:             d.Name,
							State:            false,
							FetchedTimestamp: time.Now(),
						}
					}
					dropletStates.Mutex.Unlock()
				}
			}
			if resp.Links == nil || resp.Links.IsLastPage() {
				time.Sleep(apiPollDuration)
				break PAGE_LOOP
			}
			page, err := resp.Links.CurrentPage()
			if err != nil {
				log.Printf("%v\n", err)
				time.Sleep(apiPollDuration)
				break PAGE_LOOP
			}
			opt.Page = page + 1
		}
	}
}
