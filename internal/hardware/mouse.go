package hardware

import (
    "os"
)

type MouseData struct {
    DeltaX            int
    DeltaY            int
    LeftButtonPressed bool
    RightButtonPressed bool
}

type MouseProcessor struct {
    devicePath string
    file       *os.File
}

func NewMouseProcessor(path string) *MouseProcessor {
    return &MouseProcessor{devicePath: path}
}

func (m *MouseProcessor) Initialize() error {
    f, err := os.Open(m.devicePath)
    if err != nil {
        return err
    }
    m.file = f
    return nil
}

func (m *MouseProcessor) ReadMouseData() (*MouseData, error) {
    // TODO: Parse from /dev/input
    return &MouseData{}, nil
}
