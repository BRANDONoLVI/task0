cmake_minimum_required(VERSION 3.10)

project(Task0 VERSION 0.0.0)

find_package(X11 REQUIRED)
find_package(sdbus-c++ REQUIRED)

add_executable(Task0 src/main.cpp src/hardware/MouseProcessor.cpp)
target_include_directories(Task0 PUBLIC
    "${PROJECT_BINARY_DIR}"
    "${PROJECT_SOURCE_DIR}/include"
    )

target_link_libraries(Task0 PRIVATE X11::X11 SDBusCpp::sdbus-c++)