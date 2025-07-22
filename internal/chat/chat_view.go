package chat

import (
	"fmt"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alkuinvito/chat-client/internal/discovery"
	"github.com/alkuinvito/chat-client/pkg/views"
	"github.com/hashicorp/mdns"
)

func ChatView(v *views.View) fyne.CanvasObject {
	username, err := v.Store().Get("username")
	if err != nil {
		panic("no username")
	}

	go discovery.BroadcastService(username)

	// Thread-safe chat items slice
	var (
		chatItems []string
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
			o.(*widget.Label).SetText(chatItems[i])
		},
	)

	refreshDiscovery := func() {
		entries := make(chan *mdns.ServiceEntry)

		go discovery.DiscoverService(entries)

		go func() {
			var newItems []string
			for entry := range entries {
				log.Println(entry.Host, entry.AddrV4.String())
				newItems = append(newItems, fmt.Sprintf("%s - %s", entry.Host, entry.AddrV4.String()))
			}

			fyne.Do(func() {
				mu.Lock()
				chatItems = newItems
				mu.Unlock()

				chatList.Refresh()
			})
		}()
	}

	refreshBtn := widget.NewButton("Refresh", refreshDiscovery)

	chatTabs := container.NewAppTabs(
		container.NewTabItem("New Chat", refreshBtn),
	)

	split := container.NewHSplit(chatList, chatTabs)
	split.SetOffset(0.3)

	return split
}
