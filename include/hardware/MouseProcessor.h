#ifndef MOUSEPROCESSOR_H
#define MOUSEPROCESSOR_H

#include <string>
#include "common/MouseData.h"

class MouseProcessor {
public:
    explicit MouseProcessor(const std::string& devicePath);
    ~MouseProcessor();

    bool Initialize();
    bool ReadMouseData(MouseData& data);
    void Shutdown();

    bool IsInitialized() const;

private:
    std::string devicePath;
    int deviceFd;
    bool initialized;
    bool leftButtonDown;
    bool rightButtonDown;

    void ProcessEvent(const struct input_event& ev, MouseData& data);
};

#endif