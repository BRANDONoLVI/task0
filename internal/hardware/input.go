package hardware

import (
    "fmt"
    "os"
)

func ReadMouseEvents() {
    file, err := os.Open("/dev/input/mice") // Open mouse input device
    if err != nil {
        fmt.Println("Error opening device:", err)
        return
    }
    defer file.Close()

    fmt.Println("Listening for mouse events...")
    
    buf := make([]byte, 3) // Mouse events are usually 3 bytes
    for {
        _, err := file.Read(buf)
        if err != nil {
            fmt.Println("Error reading mouse events:", err)
            continue
        }

        leftButton := buf[0] & 0x1
        rightButton := buf[0] & 0x2
        middleButton := buf[0] & 0x4
        xMove := int8(buf[1]) // X-axis movement
        yMove := int8(buf[2]) // Y-axis movement

        fmt.Printf("Mouse Event: X=%d Y=%d Left=%d Right=%d Middle=%d\n", xMove, yMove, leftButton, rightButton, middleButton)
    }
}
