#ifndef MOUSE_SERVICE_H
#define MOUSE_SERVICE_H

#include <sdbus-c++/sdbus-c++.h>
#include "hardware/MouseProcessor.h"

class MouseService {
    public:
        MouseService(*MouseProcessor);
        void Run();

    private:
        void RegisterInterface();

};

#endif