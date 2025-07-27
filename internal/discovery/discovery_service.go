package discovery

import (
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"log"
	"strings"
	"sync"
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
	s   *store.Store
}

type IDiscoveryService interface {
	BroadcastService(username string)
	GetPeers() response.Response[[]PeerModel]
	getTxt(entry *mdns.ServiceEntry, key string) string
	QueryService(entries chan *mdns.ServiceEntry)
	Startup(ctx context.Context)
}

func NewDiscoveryService(s *store.Store) *DiscoveryService {
	return &DiscoveryService{s: s}
}

func (ds *DiscoveryService) BroadcastService(id, username string) {
	txt := []string{"ID=" + id, "USERNAME=" + username}
	service, _ := mdns.NewMDNSService(id, SVC_NAME, SVC_DOMAIN, "", SVC_PORT, nil, txt)

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	<-ds.ctx.Done()
}

func (ds *DiscoveryService) GetPeers() response.Response[[]PeerModel] {
	entries := make(chan *mdns.ServiceEntry, 4)
	var result []PeerModel
	var wg sync.WaitGroup

	username, err := ds.s.GetString("user:username")
	if err != nil {
		return response.New(result).Status(500)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entries {
			id := ds.getTxt(entry, "ID")
			peerName := ds.getTxt(entry, "USERNAME")

			if peerName == username {
				continue
			}

			if id != "" && peerName != "" {
				result = append(result, PeerModel{ID: id, Username: peerName, IP: entry.AddrV4.String()})
			}
		}
	}()

	go ds.QueryService(entries)

	wg.Wait()

	return response.New(result)
}

func (ds *DiscoveryService) getTxt(entry *mdns.ServiceEntry, key string) string {
	fields := entry.InfoFields
	for _, field := range fields {
		info := strings.Split(field, "=")
		if len(info) == 2 {
			if info[0] == key {
				return info[1]
			}
		}
	}

	return ""
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
