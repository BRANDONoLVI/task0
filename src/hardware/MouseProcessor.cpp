#include "hardware/MouseProcessor.h"
#include <fcntl.h>
#include <unistd.h>
#include <linux/input.h>
#include <iostream>
#include <cstring>

MouseProcessor::MouseProcessor(const std::string& devicePath)
    : devicePath(devicePath), deviceFd(-1), initialized(false),
      leftButtonDown(false), rightButtonDown(false) {}

MouseProcessor::~MouseProcessor() {
    if (this->initialized) {
        Shutdown();
    }
}

bool MouseProcessor::Initialize() {
    this->deviceFd = open(devicePath.c_str(), O_RDONLY | O_NONBLOCK);
    if (deviceFd < 0) {
        std::cerr << "Failed to open device: " << devicePath << std::endl;
        return false;
    }
    this->initialized = true;
    return true;
}

void MouseProcessor::Shutdown() {
    if (this->deviceFd >= 0) {
        close(this->deviceFd);
        this->deviceFd = -1;
    }
    this->initialized = false;
}

bool MouseProcessor::IsInitialized() const {
    return this->initialized;
}

bool MouseProcessor::ReadMouseData(MouseData& data) {
    if (!this->initialized) {
        std::cerr << "MouseProcessor is not initialized!" << std::endl;
        return false;
    }

    struct input_event ev;

    memset(&ev, 0, sizeof(ev));
    memset(&data, 0, sizeof(MouseData));

    ssize_t bytesRead = read(this->deviceFd, &ev, sizeof(ev));

    if (bytesRead > 0) {
        ProcessEvent(ev, data);
    }

    return true;
}

void MouseProcessor::ProcessEvent(const struct input_event& ev, MouseData& data) {
    if (ev.type == EV_REL) {
        if (ev.code == REL_X) {
            data.deltaX += ev.value;
        } else if (ev.code == REL_Y) {
            data.deltaY += ev.value;
        }
    } else if (ev.type == EV_KEY) {
        if (ev.code == BTN_LEFT) {
            leftButtonDown = (ev.value != 0);
            data.leftButtonPressed = leftButtonDown;
        } else if (ev.code == BTN_RIGHT) {
            rightButtonDown = (ev.value != 0);
            data.rightButtonPressed = rightButtonDown;
        }
    }
}