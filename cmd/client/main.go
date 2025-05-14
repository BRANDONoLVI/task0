// cmd/client/main.go
package main

import (
	"fmt"
	"log"
	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	obj := conn.Object("org.luxvitae.ActionService", "/org/luxvitae/ActionService")

	gesture := "click" // or get from user input or gesture recognizer
	fmt.Println("Sending gesture:", gesture)

	call := obj.Call("org.luxvitae.ActionService.ReceiveGesture", 0, gesture)
	if call.Err != nil {
		log.Fatalf("Failed to call ReceiveGesture: %v", call.Err)
	}

	fmt.Println("Gesture sent successfully")
}