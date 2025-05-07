package processing

import (
    "time"
)

// Gesture represents a standardized user action
type Gesture struct {
    Timestamp int64
    DeviceID  string
    Type      string // e.g., "click", "drag", "move"
    Position  struct {
        X int8
        Y int8
    }
}

// DetectGesture interprets raw mouse movement
func DetectGesture(xMove, yMove int8, leftButtonPressed bool) Gesture {
    var gestureType string

    if leftButtonPressed && (xMove != 0 || yMove != 0) {
        gestureType = "drag"
    } else if leftButtonPressed {
        gestureType = "click"
    } else if xMove != 0 || yMove != 0 {
        gestureType = "move"
    }

    return Gesture{
        Timestamp: time.Now().Unix(),
        DeviceID:  "mouse_001",
        Type:      gestureType,
        Position: struct {
            X int8
            Y int8
        }{X: xMove, Y: yMove},
    }
}
