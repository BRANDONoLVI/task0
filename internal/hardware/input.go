package hardware

import (
    "fmt"
    "os"
    "time"
    "mouse-service/v1/pkg/dbus"
)

func ReadMouseEvents() {
    file, err := os.Open("/dev/input/mice") // Open mouse input device
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
        rightButton := buf[0] & 0x2
        middleButton := buf[0] & 0x4
        xMove := int8(buf[1])
        yMove := int8(buf[2])

        fmt.Printf("Mouse Event: X=%d Y=%d Left=%d Right=%d Middle=%d\n", xMove, yMove, leftButton, rightButton, middleButton)

        // Send event over D-Bus
        dbusService.SendMouseEvent("mouse_001", int(xMove), int(yMove), "pressed")
        time.Sleep(50 * time.Millisecond) // Avoid sending excessive events
    }
}
