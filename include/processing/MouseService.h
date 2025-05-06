#ifndef MOUSE_SERVICE_H
#define MOUSE_SERVICE_H

#include <sdbus-c++/sdbus-c++.h>
#include "hardware/MouseProcessor.h"

class MouseService {
    public:
        MouseService(std::shared_ptr<MouseProcessor> p);
        ~MouseService();
        void Run();

    private:
        void RegisterInterface();
        void MouseService::GetMouseState(sdbus::MethodCall& methodCall);

        std::shared_ptr<sdbus::IConnection> connection;
        std::unique_ptr<sdbus::IObject> object;
        std::shared_ptr<MouseProcessor> processor;
        sdbus::ObjectPath* path;

};

#endif