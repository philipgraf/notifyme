package dbus

import (
	"html"
	"log"

	"github.com/guelfey/go.dbus"
)

type Notification struct {
	App        string
	ReplacesId uint32
	AppIcon    string
	Summary    string
	Body       string
	Actions    []string
	Hints      map[string]dbus.Variant
	Timeout    int32
}

func (noti *Notification) String() string {
	return noti.Summary + " " + noti.Body
}

func FromMessage(msg *dbus.Message) *Notification {
	noti := &Notification{}
	noti.App = msg.Body[0].(string)
	noti.ReplacesId = msg.Body[1].(uint32)
	noti.AppIcon = msg.Body[2].(string)
	noti.Summary = msg.Body[3].(string)
	noti.Body = html.UnescapeString(msg.Body[4].(string))
	noti.Actions = msg.Body[5].([]string)
	noti.Hints = msg.Body[6].(map[string]dbus.Variant)
	noti.Timeout = msg.Body[7].(int32)
	return noti
}

func ListenForNotifications(channel chan *Notification) {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"eavesdrop='true',type='method_call',interface='org.freedesktop.Notifications',member='Notify'")
	if call.Err != nil {
		log.Fatal(call.Err)
	}
	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)
	for message := range c {
		noti := FromMessage(message)
		channel <- noti
	}
}
