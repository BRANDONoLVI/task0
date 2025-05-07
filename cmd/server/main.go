package main

import (
    "os"
    "os/signal"
    "fmt"
    "syscall"
	"mouse-service/v1/pkg/dbus"
)

func main() {
	fmt.Println("Registering D-Bus service...")
	conn, err := dbus.RegisterActionService()
    if err != nil {
        fmt.Println("D-Bus registration failed:", err)
        return
    }
	defer conn.Close()

    fmt.Println("Service running. Press Ctrl+C to exit.")
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    fmt.Println("Shutting down...")

    select {}
}