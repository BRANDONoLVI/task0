package main

import (
    "log"
    "mouse-service/v1/internal/dbus"
    "mouse-service/v1/internal/hardware"
)

func main() {
	// cat /proc/bus/input/devices
    processor := hardware.NewMouseProcessor("/dev/input/event5") // Replace with correct device
    if err := processor.Initialize(); err != nil {
        log.Fatal("Failed to initialize mouse processor:", err)
    }

    service := dbus.NewMouseService(processor)
    service.Run()
}
