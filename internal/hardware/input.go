package hardware

import (
    "fmt"
    "os"
    "time"
    "mouse-service/v1/pkg/dbus"
    "mouse-service/v1/internal/processing"
)

func ReadMouseEvents() {
    file, err := os.Open("/dev/input/mice")
    if err != nil {
        fmt.Println("Error opening device:", err)
        return
    }
    defer file.Close()

    fmt.Println("Listening for mouse events...")

    dbusService, err := dbus.NewMouseService()
    if err != nil {
        fmt.Println("Failed to initialize D-Bus service:", err)
        return
    }

    buf := make([]byte, 3)
    for {
        _, err := file.Read(buf)
        if err != nil {
            fmt.Println("Error reading mouse events:", err)
            continue
        }

        leftButton := buf[0] & 0x1
        rightButton := (buf[0] & 0x2) >> 1 // extract the second bit and shift it to the least significant bit
        // middleButton := (buf[0] & 0x4) >> 2
        xMove := int8(buf[1])
        yMove := int8(buf[2])

        gesture := processing.DetectGesture(xMove, yMove, leftButton == 1, rightButton == 1)

        //fmt.Printf("Gesture Detected: %s at (%d, %d)\n", gesture.Type, gesture.Position.X, gesture.Position.Y)

        dbusService.SendMouseEvent(gesture.DeviceID, int32(gesture.Position.X), int32(gesture.Position.Y), gesture.Type)
        //dbusService.SendMouseEvent("mouse_001", int32(xMove), int32(yMove), "pressed")
        time.Sleep(50 * time.Millisecond) // avoid spammy events
    }
}
