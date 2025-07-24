package chat

import "net"

type ChatRoom struct {
	PeerName string
	IP       net.IP
}
