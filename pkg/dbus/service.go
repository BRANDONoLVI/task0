package dbus

import (
    "fmt"
    "github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

type MouseService struct {
    conn *dbus.Conn
}

func NewMouseService() (*MouseService, error) {
    conn, err := dbus.SessionBus()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to D-Bus: %v", err)
    }

    return &MouseService{conn: conn}, nil
}
	
func (m *MouseService) SendMouseEvent(deviceID string, x, y int32, button string) {
    obj := m.conn.Object("com.example.MouseService", "/com/example/MouseService")
    call := obj.Call("com.example.MouseService.SendEvent", 0, deviceID, x, y, button)
    
    if call.Err != nil {
        fmt.Println("Failed to send event:", call.Err)
    }
}

func (m *MouseService) SendEvent(deviceID string, x int, y int, button string) *dbus.Error {
    fmt.Printf("Mouse Event Received: Device=%s X=%d Y=%d Button=%s\n", deviceID, x, y, button)
    return nil
}

func RegisterMouseService() (*dbus.Conn, error) {
    conn, err := dbus.SessionBus()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to D-Bus: %v", err)
    }

    _, err = conn.RequestName("com.example.MouseService", dbus.NameFlagDoNotQueue)
    if err != nil {
        return nil, fmt.Errorf("failed to request D-Bus name: %v", err)
    }

    // Create the service object
    service := &MouseService{}

    // Register the interface with methods
    conn.Export(service, "/com/example/MouseService", "com.example.MouseService")

    // Introspection support for debugging
    conn.Export(introspect.Introspectable("/com/example/MouseService"), "/com/example/MouseService", "org.freedesktop.DBus.Introspectable")

    fmt.Println("D-Bus service registered successfully!")

    return conn, nil
}