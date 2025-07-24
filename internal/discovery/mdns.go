package discovery

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/mdns"
)

const (
	SVC_INSTANCE = "p2p-chat"
	SVC_NAME     = "_p2pchat._tcp"
	SVC_DOMAIN   = "local."
	SVC_PORT     = 60606
)

func BroadcastService(username string) {
	service, _ := mdns.NewMDNSService(username, SVC_NAME, SVC_DOMAIN, "", SVC_PORT, nil, nil)

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func DiscoverService(entries chan *mdns.ServiceEntry) {
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
