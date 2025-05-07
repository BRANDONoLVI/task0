package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"log"
)

type MouseService struct{}

func (m *MouseService) GetMouseState() (int32, int32, bool, bool, *dbus.Error) {
	// stub values for now
	return 10, 5, true, false, nil
}

const (
	serviceName = "org.example.Mouse"
	objectPath  = "/org/example/Mouse"
	ifaceName   = "org.example.Mouse"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	reply, err := conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalf("Failed to request name: %v", err)
	}

	mouse := &MouseService{}
	conn.Export(mouse, objectPath, ifaceName)

	node := &introspect.Node{
		Name: "/",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name: ifaceName,
				Methods: []introspect.Method{
					{
						Name: "GetMouseState",
						Args: []introspect.Arg{
							{Name: "dx", Type: "i"},
							{Name: "dy", Type: "i"},
							{Name: "left", Type: "b"},
							{Name: "right", Type: "b"},
						},
					},
				},
			},
		},
	}
	conn.Export(introspect.NewIntrospectable(node), objectPath, "org.freedesktop.DBus.Introspectable")

	fmt.Println("Mouse service is running.")
	select {} // run forever
}