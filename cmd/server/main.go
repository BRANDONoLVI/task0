package main

import (
    "fmt"
    "mouse-service/v1/internal/hardware"
)

func main() {
    fmt.Println("Starting mouse event listener...")
    hardware.ReadMouseEvents()
}