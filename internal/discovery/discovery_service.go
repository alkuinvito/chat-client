package discovery

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/mdns"
)

const (
	SVC_NAME   = "_p2pchat._tcp"
	SVC_DOMAIN = "local."
	SVC_PORT   = 60606
)

type DiscoveryService struct {
	ctx context.Context
}

type IDiscoveryService interface {
	BroadcastService(username string)
	QueryService(entries chan *mdns.ServiceEntry)
	Startup(ctx context.Context)
}

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{}
}

func (ds *DiscoveryService) BroadcastService(username string) {
	service, _ := mdns.NewMDNSService(username, SVC_NAME, SVC_DOMAIN, "", SVC_PORT, nil, nil)

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	<-ds.ctx.Done()
}

func (ds *DiscoveryService) QueryService(entries chan *mdns.ServiceEntry) {
	defer close(entries)

	params := &mdns.QueryParam{
		Service: SVC_NAME,
		Domain:  SVC_DOMAIN,
		Entries: entries,
		Timeout: time.Second * 5,
	}

	err := mdns.Query(params)
	if err != nil {
		log.Fatalln(err)
	}
}

func (ds *DiscoveryService) Startup(ctx context.Context) {
	ds.ctx = ctx
}
