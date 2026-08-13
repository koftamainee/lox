cmake_minimum_required(VERSION 3.16)

project(clox LANGUAGES C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

add_executable(clox src/main.c)

target_include_directories(clox PRIVATE src)

target_compile_options(clox PRIVATE -Wall -Wextra -Wpedantic)
