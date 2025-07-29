package discovery

import (
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	SVC_NAME   = "_p2pchat._tcp"
	SVC_DOMAIN = "local."
	SVC_PORT   = 60606
)

type DiscoveryService struct {
	ctx   context.Context
	s     *store.Store
	peers map[string]*PeerModel
	mu    sync.Mutex
}

type IDiscoveryService interface {
	BroadcastService(username string)
	GetPeer(peerId string) PeerModel
	GetPeers() response.Response[[]PeerModel]
	getTxt(entry *zeroconf.ServiceEntry, key string) string
	QueryService()
	RefreshQuery()
	Startup(ctx context.Context)
}

func NewDiscoveryService(s *store.Store) *DiscoveryService {
	peers := make(map[string]*PeerModel)

	return &DiscoveryService{s: s, peers: peers}
}

func (ds *DiscoveryService) BroadcastService(id, username string) {
	txt := []string{"ID=" + id, "USERNAME=" + username}
	server, err := zeroconf.Register(id, SVC_NAME, SVC_DOMAIN, SVC_PORT, txt, nil)
	if err != nil {
		log.Println(err)
	}
	defer server.Shutdown()

	<-ds.ctx.Done()
	log.Println("Shutting down service broadcast...")
}

func (ds *DiscoveryService) GetPeer(peerId string) PeerModel {
	var result PeerModel

	ds.mu.Lock()
	for peer := range ds.peers {
		if peer == peerId {
			result = *ds.peers[peer]
		}
	}
	ds.mu.Unlock()

	return result
}

func (ds *DiscoveryService) GetPeers() response.Response[[]PeerModel] {
	var result []PeerModel

	ds.mu.Lock()
	for peer := range ds.peers {
		result = append(result, *ds.peers[peer])
	}
	ds.mu.Unlock()

	return response.New(result)
}

func (ds *DiscoveryService) getTxt(entry *zeroconf.ServiceEntry, key string) string {
	fields := entry.Text
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

func (ds *DiscoveryService) QueryService() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err)
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		username, err := ds.s.GetString("user:username")
		if err != nil {
			return
		}

		for {
			select {
			case entry := <-entries:
				if len(entry.AddrIPv4) == 0 {
					continue
				}

				peerId := ds.getTxt(entry, "ID")
				peerName := ds.getTxt(entry, "USERNAME")

				if peerName == username {
					continue
				}

				if peerId != "" && peerName != "" {
					ip := entry.AddrIPv4[0].String()
					ds.mu.Lock()
					ds.peers[entry.Instance] = &PeerModel{
						ID:       peerId,
						Username: peerName,
						IP:       ip,
					}
					ds.mu.Unlock()
				}
			case <-ds.ctx.Done():
				log.Println("Shutting down mDNS watcher...")
				return
			}
		}
	}()

	err = resolver.Browse(ds.ctx, SVC_NAME, SVC_DOMAIN, entries)
	if err != nil {
		log.Println("Failed to start browse:", err.Error())
	}
}

func (ds *DiscoveryService) RefreshQuery() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err)
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func(results chan *zeroconf.ServiceEntry) {
		username, err := ds.s.GetString("user:username")
		if err != nil {
			return
		}

		for entry := range results {
			if len(entry.AddrIPv4) == 0 {
				continue
			}

			peerId := ds.getTxt(entry, "ID")
			peerName := ds.getTxt(entry, "USERNAME")

			if peerName == username {
				continue
			}

			if peerId != "" && peerName != "" {
				ip := entry.AddrIPv4[0].String()
				ds.mu.Lock()
				ds.peers[entry.Instance] = &PeerModel{
					ID:       peerId,
					Username: peerName,
					IP:       ip,
				}
				ds.mu.Unlock()
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(ds.ctx, time.Second*5)
	defer cancel()
	err = resolver.Browse(ctx, SVC_NAME, SVC_DOMAIN, entries)
	if err != nil {
		log.Println("Failed to start browse:", err.Error())
	}

	<-ctx.Done()
}

func (ds *DiscoveryService) Startup(ctx context.Context) {
	ds.ctx = ctx
}
