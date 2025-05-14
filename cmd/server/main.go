package main

import (
    "os"
    "os/signal"
    "fmt"
    "syscall"
	"mouse-service/v1/pkg/dbus"
    "mouse-service/v1/internal/hardware"
)

func main() {
	fmt.Println("Registering D-Bus service...")
    mousConn, err := dbus.RegisterMouseService()
    if err != nil {
        fmt.Println("D-Bus Mouse registration failed:", err)
        return
    }
	defer mousConn.Close()

	conn, err := dbus.RegisterActionService()
    if err != nil {
        fmt.Println("D-Bus registration failed:", err)
        return
    }
	defer conn.Close()

    // Read mouse inputs here
    go func() {
        fmt.Println("Starting to read mouse events...")
        hardware.ReadMouseEvents()
    }()
    

    fmt.Println("Service running. Press Ctrl+C to exit.")
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    fmt.Println("Shutting down...")

    select {}
}