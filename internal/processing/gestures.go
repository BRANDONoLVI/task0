package processing

import (
    "time"
)

// this is my standardized interface
type Gesture struct {
    Timestamp int64
    DeviceID  string
    Type      string // e.g., click, drag, move, 
    Position  struct {
        X int8
        Y int8
    }
    Button    string // left, right, middle
    Direction string // optional: left, right, up, down
}

var lastClickTime int64
var lastClickPos struct{ X, Y int8 }

const SWIPE_MAGNITUDE int8 = 50 // just an arbitrary value

// Interprets raw mouse data
func DetectGesture(xMove, yMove int8, leftButtonPressed, rightButtonPressed bool) Gesture {    
    now := time.Now().UnixMilli()
    var gestureType, button, direction string
    
    if leftButtonPressed {
        button = "left"
        if xMove != 0 || yMove != 0 {
            gestureType = "drag"
        } else {
            // check for double click
            if now - lastClickTime < 500 && lastClickPos.X == xMove && lastClickPos.Y == yMove {
                gestureType = "double_click"
            } else {
                gestureType = "click"
            }
            lastClickTime = now
            lastClickPos = struct{ X, Y int8 }{xMove, yMove}
        }
    } else if rightButtonPressed {
        button = "right"
        gestureType = "right_click"
    } else if xMove != 0 || yMove != 0 {
        gestureType = "move"
        if abs(xMove) > abs(yMove) {
            if xMove > SWIPE_MAGNITUDE {
                gestureType = "swipe_right"
                direction = "right"
            } else if xMove < -SWIPE_MAGNITUDE {
                gestureType = "swipe_left"
                direction = "left"
            }
        } else {
            if yMove > SWIPE_MAGNITUDE {
                gestureType = "swipe_down"
                direction = "down"
            } else if yMove < -SWIPE_MAGNITUDE {
                gestureType = "swipe_up"
                direction = "up"
            }
        }
    }

    return Gesture{
        Timestamp: now,
        DeviceID:  "mouse_001",
        Type:      gestureType,
        Button:    button,
        Direction: direction,
        Position: struct {
            X int8
            Y int8
        }{X: xMove, Y: yMove},
    }
}

func abs(x int8) int8 {
    if x < 0 {
        return -x
    }
    return x
}