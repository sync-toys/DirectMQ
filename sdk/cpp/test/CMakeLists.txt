cmake_minimum_required(VERSION 3.26)

project(directmq-sdk-test LANGUAGES CXX)

include(FetchContent)
include(CTest)

FetchContent_Declare(
  catch2
  GIT_REPOSITORY https://github.com/catchorg/Catch2.git
  GIT_TAG v2.13.7  # Replace with the desired version
)
FetchContent_MakeAvailable(catch2)
include(${catch2_SOURCE_DIR}/contrib/Catch.cmake)

add_executable(directmq-sdk-test
  main.cpp
  test.cpp
)

target_compile_definitions(directmq-sdk-test PRIVATE PB_ENABLE_MALLOC=1)

target_compile_options(directmq-sdk-test
  PRIVATE -Wall -Wextra -Werror)

target_link_libraries(directmq-sdk-test PRIVATE
    protobuf-nanopb-static
    dmq_proto_embedded
    directmq-sdk
    Catch2::Catch2)

catch_discover_tests(directmq-sdk-test)
