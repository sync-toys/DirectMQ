cmake_minimum_required(VERSION 3.26)

project(directmq-sdk
  VERSION 1.0
  LANGUAGES CXX)

add_subdirectory(lib)

add_library(directmq-sdk STATIC src/main.cpp)
target_compile_definitions(protobuf-nanopb-static PUBLIC PB_ENABLE_MALLOC=1)

target_include_directories(directmq-sdk PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/include)

target_link_libraries(directmq-sdk PRIVATE
    protobuf-nanopb-static
    dmq_proto_embedded
    websockets)

option(TESTING "Enable tests." ON)

if(TESTING)
    enable_testing()
    add_subdirectory(test)
endif()

install(TARGETS directmq-sdk DESTINATION bin)
