include(CTest)

find_package(Catch2 REQUIRED)
include(Catch)

add_executable(directmq-sdk-test
  main.cpp
  test.cpp
)

target_compile_definitions(directmq-sdk-test PRIVATE PB_ENABLE_MALLOC=1)

target_compile_options(directmq-sdk-test
  PRIVATE -Wall -Wextra -Werror)

target_link_libraries(directmq-sdk-test PRIVATE
    protobuf-nanopb-static
    dmq-proto-embedded
    directmq-sdk
    Catch2::Catch2)

catch_discover_tests(directmq-sdk-test)
