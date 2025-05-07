package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "github.com/godbus/dbus/v5"
    "github.com/godbus/dbus/v5/introspect"
)

// Simple action execution for testing
func MapGestureToAction(gesture string) (string, error) {
    fmt.Println("Mapping gesture:", gesture)
    return "action_" + gesture, nil
}

func ExecuteAction(action string) error {
    fmt.Println("Executing action:", action)
    return nil
}

// ActionService represents a D-Bus service for handling gestures
type ActionService struct{}

// ReceiveGesture is the D-Bus method for handling gestures
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

func main() {
    fmt.Println("Starting ActionService debug version...")

    // Connect to the session bus
    conn, err := dbus.SessionBus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to session bus: %v\n", err)
        os.Exit(1)
    }

    // Request service name
    fmt.Println("Requesting D-Bus name: com.example.ActionService")
    reply, err := conn.RequestName("com.example.ActionService", dbus.NameFlagDoNotQueue)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to request D-Bus name: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("RequestName reply: %v (1=Primary owner, 2=In queue, 3=Exists, 4=Already owner)\n", reply)

    // Create service
    service := &ActionService{}
    
    // Export the service object
    fmt.Println("Exporting ActionService to /com/example/ActionService")
    err = conn.Export(service, "/com/example/ActionService", "com.example.ActionService")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to export ActionService: %v\n", err)
        os.Exit(1)
    }
    
    // Create introspection data with detailed method information
    fmt.Println("Creating introspection data...")
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
    fmt.Println("Exporting introspection interface...")
    err = conn.Export(introspect.NewIntrospectable(introspectData), "/com/example/ActionService", "org.freedesktop.DBus.Introspectable")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to export introspection: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Action Service registered on D-Bus successfully!")
    fmt.Println("Object path: /com/example/ActionService")
    fmt.Println("Interface: com.example.ActionService")
    fmt.Println("Method: ReceiveGesture(string)")
    fmt.Println("Waiting for D-Bus calls. Press Ctrl+C to exit.")
    
    // Test the introspection
    fmt.Println("\nTesting introspection...")
    obj := conn.Object("com.example.ActionService", "/com/example/ActionService")
    call := obj.Call("org.freedesktop.DBus.Introspectable.Introspect", 0)
    if call.Err != nil {
        fmt.Println("Failed to get introspection data:", call.Err)
    } else {
        var xmlData string
        call.Store(&xmlData)
        fmt.Println("Introspection data available.")
        fmt.Println("Sample of introspection XML:", xmlData[:100], "...")
    }
    
    fmt.Println("\nTo test, run this command in another terminal:")
    fmt.Println("dbus-send --session --type=method_call --dest=com.example.ActionService" +
               " /com/example/ActionService com.example.ActionService.ReceiveGesture string:\"swipe_up\"")

    // Handle signals for clean shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    fmt.Println("Shutting down...")
}