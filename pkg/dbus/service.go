package dbus

import (
    "fmt"
    "time"
    "github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
    "mouse-service/v1/internal/action"
)

type MouseService struct {
    conn *dbus.Conn
}

type ActionService struct{}

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

    service := &MouseService{}
    conn.Export(service, "/com/example/MouseService", "com.example.MouseService")
    conn.Export(introspect.Introspectable("/com/example/MouseService"), "/com/example/MouseService", "org.freedesktop.DBus.Introspectable")

    fmt.Println("D-Bus service registered successfully!")

    return conn, nil
}

// Actions

func (a *ActionService) ReceiveGesture(gesture string) *dbus.Error {
    fmt.Printf("Received gesture: %s\n", gesture)

    mappedAction, err := action.MapGestureToAction(gesture)
    if err != nil {
        fmt.Println("Error mapping gesture:", err)
        return nil
    }

    err = action.ExecuteAction(mappedAction)
    if err != nil {
        fmt.Println("Error executing action:", err)
    }

    return nil
}

func RegisterActionService() (*dbus.Conn, error) {
    conn, err := dbus.SessionBus()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to D-Bus: %v", err)
    }

    _, err = conn.RequestName("com.example.ActionService", dbus.NameFlagDoNotQueue)
    if err != nil {
        return nil, fmt.Errorf("failed to request D-Bus name: %v", err)
    }

    service := &ActionService{}
    conn.Export(service, "/com/example/ActionService", "com.example.ActionService")
    
    fmt.Println("Action Service registered on D-Bus!")

    go func() {
        time.Sleep(2 * time.Second)
        service.ReceiveGesture("click")
    }()

    return conn, nil
}