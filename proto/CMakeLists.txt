cmake_minimum_required(VERSION 3.26)

project(dmq-proto-embedded
  VERSION 1.0
  LANGUAGES CXX C)

file(
    GLOB_RECURSE
    DMQ_PROTO_EMBEDDED_SOURCES
    directmq/v1/*.c
    directmq/v1/*.h
)

add_library(dmq-proto-embedded STATIC
    ${DMQ_PROTO_EMBEDDED_SOURCES}
)

target_link_libraries(dmq-proto-embedded PRIVATE
    protobuf-nanopb-static
)

target_include_directories(dmq-proto-embedded PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/directmq/v1
)
