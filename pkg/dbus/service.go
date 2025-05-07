package dbus

import (
    "fmt"
    "github.com/godbus/dbus/v5"
)

func SendMouseEvent(deviceID string, x, y int, button string) {
    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Println("Failed to connect to D-Bus:", err)
        return
    }
    
    obj := conn.Object("com.example.MouseService", "/com/example/MouseService")
    call := obj.Call("com.example.MouseService.SendEvent", 0, deviceID, x, y, button)
    
    if call.Err != nil {
        fmt.Println("Failed to send event:", call.Err)
    }
}