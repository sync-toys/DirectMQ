cmake_minimum_required(VERSION 3.26)

include(FetchContent)

project(directmq-sdk-agent
  VERSION 1.0
  LANGUAGES CXX)

# sdk is symbolic link to sdk/cpp
add_subdirectory(sdk)

add_executable(directmq-sdk-agent main.cpp)

FetchContent_Declare(json URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz)
FetchContent_MakeAvailable(json)

target_include_directories(directmq-sdk-agent PRIVATE ${json_SOURCE_DIR}/include)

target_link_libraries(directmq-sdk-agent PRIVATE
    directmq-sdk
    protobuf-nanopb-static
    dmq_proto_embedded
    websockets
    nlohmann_json::nlohmann_json)

install(TARGETS directmq-sdk-agent DESTINATION bin)
