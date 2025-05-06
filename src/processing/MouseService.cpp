// MouseService.cpp
#include "processing/MouseService.h"
#include <iostream>

MouseService::MouseService(std::shared_ptr<MouseProcessor> pProcessor)
    : processor(pProcessor) {
    this->connection = sdbus::createSystemBusConnection();
    this->path = new sdbus::ObjectPath("/org/example/Mouse");
    this->object = sdbus::createObject(*this->connection, *path);
    RegisterInterface();
}

void MouseService::RegisterInterface() {
    const char* interfaceNameStr = "org.example.Mouse";
    sdbus::InterfaceName interfaceName{interfaceNameStr};
    sdbus::IObject* object = this->object.get();
    object->addVTable(
        sdbus::MethodVTableItem{sdbus::MethodName{"GetMouseState"}, 
        sdbus::Signature{}, 
        {}, 
        sdbus::Signature{"iiiib"}, 
        {}, 
    }).forInterface(interfaceName);
}

void MouseService::Run() {
    this->connection->enterEventLoop();
}

MouseService::~MouseService() {
    delete this->path;
    delete this->object.get();
    delete this->connection.get();
}