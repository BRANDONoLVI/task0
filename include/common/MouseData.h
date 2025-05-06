#ifndef MOUSEDATA_H
#define MOUSEDATA_H

struct MouseData {
    int deltaX;
    int deltaY;
    bool leftButtonPressed;
    bool rightButtonPressed;

    MouseData() = default;
};

#endif