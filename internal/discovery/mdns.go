package discovery

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/mdns"
)

const (
	SVC_INSTANCE = "p2p-chat"
	SVC_TYPE     = "_p2pchat._tcp"
	SVC_DOMAIN   = "local."
	SVC_PORT     = 60606
)

func BroadcastService(username string) {
	service, _ := mdns.NewMDNSService(username, SVC_TYPE, "", "", SVC_PORT, nil, nil)

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func DiscoverService(entries chan *mdns.ServiceEntry) {
	err := mdns.Lookup(SVC_TYPE, entries)
	if err != nil {
		log.Fatalln(err)
	}

	close(entries)
}
