cmake_minimum_required(VERSION 3.26)

include(FetchContent)

# Fetch nanopb
FetchContent_Declare(
  nanopb
  GIT_REPOSITORY https://github.com/nanopb/nanopb.git
  GIT_TAG 0.4.9
)
FetchContent_MakeAvailable(nanopb)

add_subdirectory(dmq_proto_embedded)
