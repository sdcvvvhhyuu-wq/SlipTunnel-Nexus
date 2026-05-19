cmake_minimum_required(VERSION 3.22.1)
project("argotunnel_jni")

add_library(argotunnel_jni SHARED native-lib.cpp)

# Link against the Go C-Shared library generated in Phase 1
add_library(argotunnel_core SHARED IMPORTED)
set_target_properties(argotunnel_core PROPERTIES IMPORTED_LOCATION
    ${CMAKE_SOURCE_DIR}/../jniLibs/${ANDROID_ABI}/libargotunnel_core.so)

# THE FIX: Tell CMake where to find libargotunnel_core.h
target_include_directories(argotunnel_jni PRIVATE
    ${CMAKE_SOURCE_DIR}/../jniLibs/${ANDROID_ABI}
)

find_library(log-lib log)

target_link_libraries(argotunnel_jni argotunnel_core ${log-lib})
