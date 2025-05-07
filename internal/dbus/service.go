package dbus

import (
    "mouse-service/v1/internal/hardware"
)

type MouseService struct {
    processor *hardware.MouseProcessor
}

func NewMouseService(p *hardware.MouseProcessor) *MouseService {
    return &MouseService{processor: p}
}

func (s *MouseService) Run() {
    // TODO: Connect to system/user bus, export method like GetMouseState
}
