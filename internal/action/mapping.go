package action

import (
    "fmt"
    "os/exec"
)

type Action struct {
    Gesture string
    Command string
}

func MapGestureToAction(gesture string) (Action, error) {
    fmt.Printf("Found an action: %s\n", gesture)

    actionMap := map[string]string{
        "click": "xdotool click 1", 
        "drag":  "xdotool mousedown 1", 
        "move":  "xdotool mousemove_relative -- 10 0",
    }

    cmd, exists := actionMap[gesture]
    if !exists {
        return Action{}, fmt.Errorf("no action mapped for gesture: %s", gesture)
    }

    return Action{Gesture: gesture, Command: cmd}, nil
}

func ExecuteAction(action Action) error {
    fmt.Printf("Executing action: %s\n", action.Command)
    cmd := exec.Command("sh", "-c", action.Command)
    return cmd.Run()
}
