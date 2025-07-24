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

	// Thread-safe chat items slice
	var (
		chatItems []ChatRoom
		mu        sync.Mutex
	)

	chatList := widget.NewList(
		func() int {
			mu.Lock()
			defer mu.Unlock()
			return len(chatItems)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			mu.Lock()
			defer mu.Unlock()
			o.(*widget.Label).SetText(fmt.Sprintf("%s - %s", chatItems[i].IP.String(), chatItems[i].PeerName))
		},
	)

	refreshDiscovery := func() {
		entries := make(chan *mdns.ServiceEntry, 4)

		go discovery.DiscoverService(entries)

		go func() {
			var newItems []ChatRoom
			for entry := range entries {
				isP2PChat := strings.HasSuffix(entry.Name, fmt.Sprintf("%s.%s", discovery.SVC_NAME, discovery.SVC_DOMAIN))
				if isP2PChat {
					peerName := strings.Split(entry.Name, ".")[0]
					isOwnService := peerName == username
					if !isOwnService {
						newItems = append(newItems, ChatRoom{PeerName: peerName, IP: entry.AddrV4})
					}
				}
			}

			fyne.Do(func() {
				mu.Lock()
				chatItems = newItems
				mu.Unlock()

				dialog.ShowInformation("Scan result", fmt.Sprintf("Available peer(s) found: %d", len(newItems)), v.Window())

				chatList.Refresh()
			})
		}()
	}

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), refreshDiscovery)

	sideBar := container.NewBorder(nil, refreshBtn, nil, nil, chatList)

	msgInput := widget.NewEntry()
	msgInput.SetPlaceHolder("Your message here...")

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
	})

	chatInput := container.NewBorder(nil, nil, nil, sendBtn, msgInput)
	blankChat := container.NewCenter(widget.NewLabel("No chat"))

	split := container.NewHSplit(sideBar, blankChat)
	split.SetOffset(0.3)

	chatList.OnSelected = func(id widget.ListItemID) {
		mu.Lock()
		currentChat = chatItems[id]
		split.Trailing = container.NewBorder(
			widget.NewLabel(fmt.Sprintf("Connected to: %s", chatItems[id].IP.String())),
			chatInput,
			nil,
			nil,
			widget.NewLabel("No chat"),
		)
		mu.Unlock()

		split.Refresh()
	}

	return split
}
