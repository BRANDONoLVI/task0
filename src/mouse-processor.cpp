#include "mouse-processor.h"
#include <fcntl.h>
#include <unistd.h>
#include <linux/input.h>
#include <iostream>
#include <cstring>


MouseProcessor::MouseProcessor(const char* path) {
    // realistically, should check for nullptr

    // copy the string into a buffer
    size_t len = strlen(path) + 1; // add 1 for the null terminator
    char* buffer = new char[len];
    mempcpy(buffer, path, len);

    this->mousePath = buffer;

    this->deviceFd = -1;
    this->isInitialized = false;
    this->absX = 0;
    this->absY = 0;
}

MouseProcessor::~MouseProcessor() {
    if (isInitialized) {
        Shutdown();
    }
}

bool MouseProcessor::Initialize() {
    deviceFd = open(this->mousePath, O_RDONLY | O_NONBLOCK);
    if (deviceFd < 0) {
        std::cerr << "Failed to open device: " << this->mousePath << deviceFd << std::endl;
        return false;
    }

    Display* display = XOpenDisplay(nullptr);
    if (display == nullptr) {
        std::cerr << "Cannot open display" << std::endl;
        return false;
    }

    this->display = display;
    this->rootWindow = DefaultRootWindow(display);

    isInitialized = true;
    return true;
}

/**
 * https://www.kernel.org/doc/html/v4.17/input/event-codes.html
 */
bool MouseProcessor::ReadMouseData(MouseData& data) {
    if (!isInitialized) {
        std::cerr << "MouseProcessor is not initialized!" << std::endl;
        return false;
    }

    struct input_event ev;
    memset(&data, 0, sizeof(MouseData)); // reset struct rather than allocate mem

    Window ret_root, ret_child;

    while (true) {
        ssize_t bytesRead = read(deviceFd, &ev, sizeof(ev));
        if (bytesRead > 0) {
            if (ev.type == EV_REL) {
                if (ev.code == REL_X) {
                    this->absX += ev.value;
                    data.deltaX += ev.value;
                } else if (ev.code == REL_Y) {
                    this->absY += ev.value;
                    data.deltaY += ev.value;
                }
                std::cout << "(" << this->absX << ", " << this->absY << ")" << std::endl;
            } else if (ev.type == EV_KEY) {
                std::cout << "Button event: code=" << ev.code << ", value=" << ev.value << std::endl;
            }
        } else if (bytesRead == 0) {
            std::cout << "No data available." << std::endl;
        }

        int root_x, root_y;
        int win_x, win_y;
        unsigned int mask;

        Window ret_root, ret_child;
        if (XQueryPointer(display, this->rootWindow, &ret_root, &ret_child,
                        &root_x, &root_y, &win_x, &win_y, &mask)) {
            std::cout << "Absolute mouse position: (" << root_x << ", " << root_y << ")" << std::endl;
        } else {
            std::cerr << "Failed to query pointer." << std::endl;
        }
    }
    

    return true;
}

std::tuple<uint32_t, uint32_t> MouseProcessor::GetAbsoluteMousePosition() const {

}

bool MouseProcessor::IsMouseLeftButtonDown() const {
    return this->mouseLeftButtonDown;
}

bool MouseProcessor::IsMouseRightButtonDown() const {
    return this->mouseRightButtonDown;
}

void MouseProcessor::Shutdown() {
    if (deviceFd >= 0) {
        close(deviceFd);
        deviceFd = -1;
    }
    isInitialized = false;
}