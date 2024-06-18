cmake_minimum_required(VERSION 3.26)

project(directmq-sdk
  VERSION 1.0
  LANGUAGES CXX)

set(NANOPB_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/lib/nanopb")
include_directories(${NANOPB_INCLUDE_DIR})

add_subdirectory(lib/nanopb)
add_subdirectory(lib/dmq-proto-embedded)

add_library(directmq-sdk STATIC src/main.cpp)
target_compile_definitions(protobuf-nanopb-static PUBLIC PB_ENABLE_MALLOC=1)

target_include_directories(directmq-sdk PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/include)

target_link_libraries(directmq-sdk PRIVATE
    protobuf-nanopb-static
    dmq-proto-embedded)

option(TESTING "Enable tests." ON)

if(TESTING)
  enable_testing()
  add_subdirectory(test)
endif()

install(TARGETS directmq-sdk DESTINATION bin)
