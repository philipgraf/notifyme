package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/philipgraf/libtmux-go"
	"github.com/philipgraf/notifyme/dbus"
)

func main() {
	c := make(chan *dbus.Notification, 10)
	go dbus.ListenForNotifications(c)
	for noti := range c {
		c := color.New(color.FgYellow).Add(color.Bold)
		c.Printf("%s ", noti.Summary)
		fmt.Println(noti.Body)
		displayInTmux(noti)
	}
}

func displayInTmux(noti *dbus.Notification) {
	clients := tmux.GetAllClients()
	option := tmux.NewDisplayOptions()
	option.Time = "2000"
	tmux.Display(noti.String(), option, clients...)
}
