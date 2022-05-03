cmake_minimum_required(VERSION 3.10)
project(network_hw C)

set(CMAKE_C_STANDARD 11)

add_executable(network_hw
        submit_20176342_hw3/client.c
        submit_20176342_hw3/server.c)
