package chat

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alkuinvito/chat-client/internal/discovery"
	"github.com/alkuinvito/chat-client/pkg/views"
	"github.com/hashicorp/mdns"
)

func ChatView(v *views.View) fyne.CanvasObject {
	currentChat := ChatRoom{}

	username, err := v.Store().Get("username")
	if err != nil {
		panic("no username")
	}

	go discovery.BroadcastService(username)

	var (
		rooms  []ChatRoom
		muRoom sync.Mutex
	)

	var (
		chats  []ChatMessage
		muChat sync.Mutex
	)

	chatList := widget.NewList(
		func() int {
			muChat.Lock()
			defer muChat.Unlock()
			return len(chats)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			muChat.Lock()
			defer muChat.Unlock()
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s", chats[i].Sender, chats[i].Message))
		},
	)

	// run chat server to handle incoming messages
	msgStream := make(chan *ChatMessage)
	go ServeChat(msgStream)

	// loop to retrieve new message
	go func() {
		var newChats []ChatMessage
		for msg := range msgStream {
			fyne.Do(func() {
				newChats = append(newChats, *msg)

				muChat.Lock()
				chats = newChats
				muChat.Unlock()

				chatList.Refresh()
			})
		}
	}()

	roomList := widget.NewList(
		func() int {
			muRoom.Lock()
			defer muRoom.Unlock()
			return len(rooms)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			muRoom.Lock()
			defer muRoom.Unlock()
			o.(*widget.Label).SetText(fmt.Sprintf("%s - %s", rooms[i].IP.String(), rooms[i].PeerName))
		},
	)

	refreshDiscovery := func() {
		entries := make(chan *mdns.ServiceEntry, 4)

		go discovery.DiscoverService(entries)

		go func() {
			var newRooms []ChatRoom
			for entry := range entries {
				isP2PChat := strings.HasSuffix(entry.Name, fmt.Sprintf("%s.%s", discovery.SVC_NAME, discovery.SVC_DOMAIN))
				if isP2PChat {
					peerName := strings.Split(entry.Name, ".")[0]
					isOwnService := peerName == username
					if !isOwnService {
						newRooms = append(newRooms, ChatRoom{PeerName: peerName, IP: entry.AddrV4})
					}
				}
			}

			fyne.Do(func() {
				muRoom.Lock()
				rooms = newRooms
				muRoom.Unlock()

				dialog.ShowInformation("Scan result", fmt.Sprintf("Available peer(s) found: %d", len(newRooms)), v.Window())

				roomList.Refresh()
			})
		}()
	}

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), refreshDiscovery)

	sideBar := container.NewBorder(nil, refreshBtn, nil, nil, roomList)

	msgInput := widget.NewEntry()
	msgInput.SetPlaceHolder("Your message here...")
	// handle submit on enter
	msgInput.OnSubmitted = func(s string) {
		message := ChatMessage{
			Sender:  username,
			Message: s,
		}

		go func() {
			_, err := SendMessage(discovery.SVC_PORT, currentChat, message)
			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(err, v.Window())
				})
			}
		}()

		// clear input after submit
		msgInput.SetText("")
	}

	sendBtn := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), func() {
		message := ChatMessage{
			Sender:  username,
			Message: msgInput.Text,
		}

		go func() {
			_, err := SendMessage(discovery.SVC_PORT, currentChat, message)
			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(err, v.Window())
				})
			}
		}()

		// clear input after submit
		msgInput.SetText("")
	})

	chatInput := container.NewBorder(nil, nil, nil, sendBtn, msgInput)
	blankChat := container.NewCenter(widget.NewLabel("No chat"))

	split := container.NewHSplit(sideBar, blankChat)
	split.SetOffset(0.3)

	roomList.OnSelected = func(id widget.ListItemID) {
		muRoom.Lock()
		currentChat = rooms[id]
		muRoom.Unlock()

		split.Trailing = container.NewBorder(
			widget.NewLabel(fmt.Sprintf("Connected to: %s - %s", currentChat.IP.String(), currentChat.PeerName)),
			chatInput,
			nil,
			nil,
			chatList,
		)

		split.Refresh()
	}

	return split
}
