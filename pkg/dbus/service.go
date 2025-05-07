package dbus

import (
	"fmt"
	"mouse-service/v1/internal/action"
	"os"

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
    obj := m.conn.Object("org.luxvitae.MouseService", "/org/luxvitae/MouseService")
    call := obj.Call("org.luxvitae.MouseService.SendEvent", 0, deviceID, x, y, button)
    
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

    _, err = conn.RequestName("org.luxvitae.MouseService", dbus.NameFlagDoNotQueue)
    if err != nil {
        return nil, fmt.Errorf("failed to request D-Bus name: %v", err)
    }

    service := &MouseService{}
    conn.Export(service, "/org/luxvitae/MouseService", "org.luxvitae.MouseService")
    conn.Export(introspect.Introspectable("/org/luxvitae/MouseService"), "/org/luxvitae/MouseService", "org.freedesktop.DBus.Introspectable")

    fmt.Println("D-Bus service registered successfully!")

    return conn, nil
}

// Actions

func (a *ActionService) ReceiveGesture(gesture string) *dbus.Error {
    fmt.Println("=== ReceiveGesture method called with gesture:", gesture)
    
    mappedAction, err := action.MapGestureToAction(gesture)

    if err != nil {
        fmt.Println("Error mapping gesture:", err)
        return dbus.MakeFailedError(err)
    }

    err = action.ExecuteAction(mappedAction)
    if err != nil {
        fmt.Println("Error executing action:", err)
        return dbus.MakeFailedError(err)
    }

    fmt.Println("Gesture handled successfully:", gesture)
    return nil
}

func RegisterActionService() (*dbus.Conn, error) {
    fmt.Println("Starting ActionService debug version...")

    conn, err := dbus.SessionBus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to session bus: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Requesting D-Bus name: org.luxvitae.ActionService")
    reply, err := conn.RequestName("org.luxvitae.ActionService", dbus.NameFlagDoNotQueue)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to request D-Bus name: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("RequestName reply: %v (1=Primary owner, 2=In queue, 3=Exists, 4=Already owner)\n", reply)

    service := &ActionService{}
    
    fmt.Println("Exporting ActionService to /org/luxvitae/ActionService")
    err = conn.Export(service, "/org/luxvitae/ActionService", "org.luxvitae.ActionService")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to export ActionService: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Creating introspection data...")
    introspectData := &introspect.Node{
        Name: "/org/luxvitae/ActionService",
        Interfaces: []introspect.Interface{
            {
                Name: "org.luxvitae.ActionService",
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
    fmt.Println("Exporting introspection interface...")
    err = conn.Export(introspect.NewIntrospectable(introspectData), "/org/luxvitae/ActionService", "org.freedesktop.DBus.Introspectable")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to export introspection: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Action Service registered on D-Bus successfully!")
    fmt.Println("Object path: /org/luxvitae/ActionService")
    fmt.Println("Interface: org.luxvitae.ActionService")
    fmt.Println("Method: ReceiveGesture(string)")
    fmt.Println("Waiting for D-Bus calls. Press Ctrl+C to exit.")
    
    fmt.Println("\nTesting introspection...")
    obj := conn.Object("org.luxvitae.ActionService", "/org/luxvitae/ActionService")
    call := obj.Call("org.freedesktop.DBus.Introspectable.Introspect", 0)
    if call.Err != nil {
        fmt.Println("Failed to get introspection data:", call.Err)
    } else {
        var xmlData string
        call.Store(&xmlData)
        fmt.Println("Introspection data available.")
    }
    
    fmt.Println("\nTo test, run this command in another terminal:")
    fmt.Println("dbus-send --session --type=method_call --dest=org.luxvitae.ActionService" +
               " /org/luxvitae/ActionService org.luxvitae.ActionService.ReceiveGesture string:\"click\"")

    return conn, nil
}