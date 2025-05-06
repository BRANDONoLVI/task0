#include "hardware/MouseProcessor.h"
#include "processing/MouseService.h"

#include "common/MouseData.h"
#include <iostream>
#include <unistd.h>

int main() {    
    std::shared_ptr<MouseProcessor> processor = std::make_shared<MouseProcessor>("/dev/input/event5");
    if (!processor->Initialize()) return EXIT_FAILURE;

    std::cout << "MouseProcessor initialized" << std::endl;

    MouseService service(processor);
    service.Run();

    /*
    MouseData data;
    
    while (true) {
        processor.ReadMouseData(data);
        if (data.deltaX || data.deltaY) {
            std::cout << "Moved: dx=" << data.deltaX << ", dy=" << data.deltaY << "\n";
        }
        if (data.leftButtonPressed) {
            std::cout << "Left button pressed\n";
        }
        if (data.rightButtonPressed) {
            std::cout << "Right button pressed\n";
        }
        usleep(10000);
    }

    processor.Shutdown();
    */
   
    return EXIT_SUCCESS;
}
