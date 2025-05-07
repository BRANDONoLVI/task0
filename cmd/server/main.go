package main

import (
    "fmt"
    "mouse-service/v1/internal/hardware"
	"mouse-service/v1/pkg/dbus"
)

func main() {
	fmt.Println("Registering D-Bus service...")
	conn, err := dbus.RegisterMouseService()
    if err != nil {
        fmt.Println("D-Bus registration failed:", err)
        return
    }
	defer conn.Close()

    fmt.Println("Starting mouse event listener...")
    hardware.ReadMouseEvents()
}