cmake_minimum_required(VERSION 3.26)

include(FetchContent)

project(directmq-sdk-agent
  VERSION 1.0
  LANGUAGES CXX)

FetchContent_Declare(
    directmq-sdk
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/../../../sdk/cpp
)
FetchContent_MakeAvailable(directmq-sdk)

find_package(nlohmann_json REQUIRED)

add_executable(directmq-sdk-agent main.cpp)

target_include_directories(directmq-sdk-agent PRIVATE ${json_SOURCE_DIR}/include)

target_link_libraries(directmq-sdk-agent PRIVATE
    directmq-sdk
    protobuf-nanopb-static
    dmq_proto_embedded
    nlohmann_json::nlohmann_json)

install(TARGETS directmq-sdk-agent DESTINATION bin)
