package dbus

import (
    "fmt"
    "github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
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
func MapGestureToAction(gesture string) (string, error) {
    fmt.Println("Mapping gesture:", gesture)
    return "action_" + gesture, nil
}

func ExecuteAction(action string) error {
    fmt.Println("Executing action:", action)
    return nil
}

func (a *ActionService) ReceiveGesture(gesture string) *dbus.Error {
    fmt.Println("=== ReceiveGesture method called with gesture:", gesture)
    
    mappedAction, err := MapGestureToAction(gesture)
    if err != nil {
        fmt.Println("Error mapping gesture:", err)
        return dbus.MakeFailedError(err)
    }

    err = ExecuteAction(mappedAction)
    if err != nil {
        fmt.Println("Error executing action:", err)
        return dbus.MakeFailedError(err)
    }

    fmt.Println("Gesture handled successfully:", gesture)
    return nil
}

func RegisterActionService() (*dbus.Conn, error) {
    conn, err := dbus.SessionBus()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to D-Bus: %v", err)
    }

    reply, err := conn.RequestName("com.example.ActionService", dbus.NameFlagDoNotQueue)
    if err != nil {
        return nil, fmt.Errorf("failed to request D-Bus name: %v", err)
    }

    fmt.Printf("Action service request name reply: %v\n", reply)

    // Create and export the service
    service := &ActionService{}
    
    // Export the service interface
    err = conn.Export(service, "/com/example/ActionService", "com.example.ActionService")
    if err != nil {
        return nil, fmt.Errorf("failed to export service: %v", err)
    }
    
    // Create introspection data
    introspectData := &introspect.Node{
        Name: "/com/example/ActionService",
        Interfaces: []introspect.Interface{
            {
                Name: "com.example.ActionService",
                Methods: []introspect.Method{
                    {
                        Name: "ReceiveGesture",
                        Args: []introspect.Arg{
                            {Name: "gesture", Type: "s", Direction: "in"},
                        },
                    },
                },
            },
            introspect.IntrospectData,
        },
    }
    
    // Export the introspection interface
    err = conn.Export(introspect.NewIntrospectable(introspectData), "/com/example/ActionService", "org.freedesktop.DBus.Introspectable")
    if err != nil {
        return nil, fmt.Errorf("failed to export introspection: %v", err)
    }
    
    fmt.Println("Action Service registered on D-Bus!")

    return conn, nil
}