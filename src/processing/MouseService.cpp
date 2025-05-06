// MouseService.cpp
#include "processing/MouseService.h"
#include <iostream>

MouseService::MouseService(std::shared_ptr<MouseProcessor> pProcessor)
    : processor(pProcessor) {
    if (!processor) {
        std::cerr << "MouseProcessor is null" << std::endl;
        throw std::runtime_error("MouseProcessor is null");
    }
    if (!processor->IsInitialized()) {
        std::cerr << "MouseProcessor is not initialized" << std::endl;
        throw std::runtime_error("MouseProcessor is not initialized");
    }
    this->connection = sdbus::createSystemBusConnection();

    std::cout << "Connecting to system bus..." << std::endl;
    if (!this->connection) {
        std::cerr << "Failed to create system bus connection" << std::endl;
        throw std::runtime_error("Failed to create system bus connection");
    }

    this->path = new sdbus::ObjectPath("/org/example/Mouse");

    if (!this->path) {
        std::cerr << "Failed to create object path" << std::endl;
        throw std::runtime_error("Failed to create object path");
    }
    std::cout << "Creating object..." << std::endl;
    this->object = sdbus::createObject(*this->connection, *this->path);

    if (!this->object) {
        std::cerr << "Failed to create object" << std::endl;
        throw std::runtime_error("Failed to create object");
    }
    std::cout << "Registering object..." << std::endl;
    RegisterInterface();
}

void MouseService::RegisterInterface() {
    const char* interfaceNameStr = "org.example.Mouse";
    sdbus::InterfaceName interfaceName{interfaceNameStr};
    sdbus::IObject* object = this->object.get();

    std::cout << "Registering interface: " << interfaceNameStr << std::endl;

    // TODO: fix this
    object->addVTable(
        sdbus::MethodVTableItem{sdbus::MethodName{"GetMouseState"}, 
        sdbus::Signature{}, 
        {}, 
        sdbus::Signature{"iiiib"}, 
        {}, 

    }).forInterface(interfaceName);

    std::cout << "Registering method: GetMouseState" << std::endl;
}



void MouseService::GetMouseState(sdbus::MethodCall& methodCall) {
    MouseData data;
    if (this->processor->ReadMouseData(data)) {
        // create a reply
        auto reply = methodCall.createReply();
        reply << data.deltaX << data.deltaY << data.leftButtonPressed << data.rightButtonPressed;
        // send the reply
        reply.send();
    } else {
        std::cerr << "Failed to read mouse data" << std::endl;
        // TODO: create an error reply
    }

    // emit the signal

}

void MouseService::Run() {
    this->connection->enterEventLoop();
}

MouseService::~MouseService() {
    delete this->path;
    delete this->object.get();
    delete this->connection.get();
}